#!/bin/bash
set -e

# This script builds the debian package.
# The VERSION argument is passed from the GitHub workflow (e.g., v0.3.0)
VERSION_TAG=${1} 
VERSION_NUMBER=${1#v} # a v prefix from the tag

if [ -z "$VERSION_TAG" ]; then
    echo "Usage: $0 <version_tag>"
    exit 1
fi

# Create the package structure using the full tag (e.g., docker-ai_v0.3.0_amd64)
DEB_DIR="docker-ai_${VERSION_TAG}_amd64"
mkdir -p "$DEB_DIR/usr/local/bin"
mkdir -p "$DEB_DIR/DEBIAN"

# Build the Go binary
echo "Building docker-ai for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$DEB_DIR/usr/local/bin/docker-ai" ./cmd/docker-ai

# Create the debian control file
cat > "$DEB_DIR/DEBIAN/control" <<EOF
Package: docker-ai
Version: $VERSION_NUMBER
Architecture: amd64
Maintainer: Ajay <ajay@example.com>
Description: An AI-powered CLI for Docker.
EOF

# Build the .deb package
echo "Building Debian package..."
dpkg-deb --build "$DEB_DIR"

echo "Package created: ${DEB_DIR}.deb" 