#!/bin/bash
# A script to auto-generate release notes from git history.

set -e

# --- Emojis ---
FEAT="✨"
FIX="🐛"
HOTFIX="🔥"
CHORE="🔨"

# Get the latest tag, or an empty string if it doesn't exist.
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

if [ -z "$LATEST_TAG" ]; then
    echo "No previous tag found, generating notes for all commits."
    # Get all commit subjects from the very first commit.
    LOG_CMD="git log --pretty=format:'%s'"
else
    echo "Generating notes from tag $LATEST_TAG to HEAD."
    # Get all commit subjects since the last tag.
    LOG_CMD="git log $LATEST_TAG..HEAD --pretty=format:'%s'"
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

# Generate notes for each section if commits for it exist.

# --- Features ---
FEAT_COMMITS=$(grep "^feat:" "$TMP_LOG" || true)
if [ -n "$FEAT_COMMITS" ]; then
    echo "### $FEAT Features"
    echo "$FEAT_COMMITS" | sed 's/feat:/-/'
    echo ""
fi

# --- Bug Fixes ---
FIX_COMMITS=$(grep "^fix:" "$TMP_LOG" || true)
if [ -n "$FIX_COMMITS" ]; then
    echo "### $FIX Bug Fixes"
    echo "$FIX_COMMITS" | sed 's/fix:/-/'
    echo ""
fi

# --- Hotfixes ---
HOTFIX_COMMITS=$(grep "^hotfix:" "$TMP_LOG" || true)
if [ -n "$HOTFIX_COMMITS" ]; then
    echo "### $HOTFIX Hotfixes"
    echo "$HOTFIX_COMMITS" | sed 's/hotfix:/-/'
    echo ""
fi

# --- Other Commits / Chores ---
# Select all commits that don't match the main types.
OTHER_COMMITS=$(grep -v -e "^feat:" -e "^fix:" -e "^hotfix:" "$TMP_LOG" || true)
if [ -n "$OTHER_COMMITS" ]; then
    echo "### $CHORE Other Changes"
    echo "$OTHER_COMMITS" | sed 's/^/* /'
    echo ""
fi 