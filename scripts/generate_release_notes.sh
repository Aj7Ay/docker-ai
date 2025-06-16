#!/bin/bash
# A script to auto-generate release notes from git history.

set -e

# --- Emojis ---
FEAT="✨"
FIX="🐛"
HOTFIX="🔥"
CHORE="🔨"

# In a GitHub Action, GITHUB_REF_NAME will be set (e.g., 'refs/tags/v0.3.0')
# We extract the tag name from it. For local runs, we just use the latest tag.
CURRENT_TAG=${GITHUB_REF_NAME##*/v}
if [ -z "$CURRENT_TAG" ]; then
    CURRENT_TAG=$(git describe --tags --abbrev=0)
fi

# Find the tag of the commit *before* the current tag. This is our "from" point.
# If it fails (e.g., this is the very first tag), PREVIOUS_TAG will be empty.
PREVIOUS_TAG=$(git describe --tags --abbrev=0 "${CURRENT_TAG}^" 2>/dev/null || echo "")

if [ -z "$PREVIOUS_TAG" ]; then
    echo "No previous tag found. Generating notes for all commits up to ${CURRENT_TAG}."
    # Get all commit messages from the very first commit up to the current tag
    LOG_CMD="git log ${CURRENT_TAG} --pretty=format:'%B'"
else
    echo "Generating notes from tag ${PREVIOUS_TAG} to ${CURRENT_TAG}."
    # Get all commit messages between the two tags.
    LOG_CMD="git log ${PREVIOUS_TAG}..${CURRENT_TAG} --pretty=format:'%B'"
fi

# Use a temporary file to store logs to avoid issues with pipes and loops.
TMP_LOG=$(mktemp)
# Ensure the temp file is cleaned up on exit.
trap 'rm -f "$TMP_LOG"' EXIT

# Execute the git log command.
eval $LOG_CMD > "$TMP_LOG"

echo ""
echo "## What's New"
echo ""

# --- Features ---
# Grep for lines that contain "feat:", case-insensitive
FEAT_COMMITS=$(grep -i "feat:" "$TMP_LOG" || true)
if [ -n "$FEAT_COMMITS" ]; then
    echo "### $FEAT Features"
    # Clean up the line, remove the prefix, trim whitespace, and prepend a dash
    echo "$FEAT_COMMITS" | sed -e 's/.*feat://i' -e 's/^[[:space:]]*//' -e 's/^/- /'
    echo ""
fi

# --- Bug Fixes ---
FIX_COMMITS=$(grep -i "fix:" "$TMP_LOG" || true)
if [ -n "$FIX_COMMITS" ]; then
    echo "### $FIX Bug Fixes"
    echo "$FIX_COMMITS" | sed -e 's/.*fix://i' -e 's/^[[:space:]]*//' -e 's/^/- /'
    echo ""
fi

# --- Hotfixes ---
HOTFIX_COMMITS=$(grep -i "hotfix:" "$TMP_LOG" || true)
if [ -n "$HOTFIX_COMMITS" ]; then
    echo "### $HOTFIX Hotfixes"
    echo "$HOTFIX_COMMITS" | sed -e 's/.*hotfix://i' -e 's/^[[:space:]]*//' -e 's/^/- /'
    echo ""
fi

# --- Other Commits / Chores ---
# Select all commits that don't match the main types.
OTHER_COMMITS=$(grep -v -e "^feat:" -e "^fix:" -e "^hotfix:" "$TMP_LOG" || true)
if [ -n "$OTHER_COMMITS" ]; then
    echo "### $CHORE Other Changes"
    # Format multi-line commit messages correctly.
    # The first line of a commit gets a bullet, subsequent lines are indented.
    echo "$OTHER_COMMITS" | awk 'BEGIN{first=1} /^$/{first=1; next} {if(first){print "* " $0; first=0} else {print "  " $0}}'
    echo ""
fi 