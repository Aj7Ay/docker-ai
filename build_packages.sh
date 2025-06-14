#!/bin/bash
set -e

# --- Configuration ---
APP_NAME="docker-ai"
VERSION="0.1.0" # IMPORTANT: Update this version for new releases
MAINTAINER="Ajay Kumar <postbox.aj99@gmail.com>" # IMPORTANT: Change this
DESCRIPTION="An AI-powered CLI for Docker."
HOMEPAGE="https://github.com/Aj7Ay/docker-ai"
ARCH=$(dpkg --print-architecture)

# --- Script ---
echo "Starting build for ${APP_NAME} v${VERSION} for ${ARCH}..."

# Create a temporary directory for packaging
PACKAGE_DIR="${APP_NAME}_${VERSION}_${ARCH}"
rm -rf "${PACKAGE_DIR}"
mkdir -p "${PACKAGE_DIR}/DEBIAN"
mkdir -p "${PACKAGE_DIR}/usr/local/bin"

# Build the Go binary for Linux
echo "Building Go binary..."
GOOS=linux GOARCH=${ARCH} go build -o "${APP_NAME}" ./cmd/docker-ai/main.go

# Move binary to the package directory
mv "${APP_NAME}" "${PACKAGE_DIR}/usr/local/bin/"

# Create the control file
echo "Creating DEBIAN/control file..."
INSTALLED_SIZE=$(du -sk "${PACKAGE_DIR}" | awk '{print $1}')

cat > "${PACKAGE_DIR}/DEBIAN/control" <<EOF
Package: ${APP_NAME}
Version: ${VERSION}
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Installed-Size: ${INSTALLED_SIZE}
Description: ${DESCRIPTION}
Homepage: ${HOMEPAGE}
Section: utils
Priority: optional
EOF

# Build the .deb package
echo "Building .deb package..."
dpkg-deb --build "${PACKAGE_DIR}"

echo "Successfully created ${PACKAGE_DIR}.deb"

# Cleanup
rm -rf "${PACKAGE_DIR}" 