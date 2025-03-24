#!/bin/bash

set -e

REPO="kyco/steuer-go"
INSTALL_DIR="/usr/local/bin"
EXECUTABLE="steuergo"

# Determine the OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture to GitHub release format
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  i386 | i686)
    ARCH="386"
    ;;
  aarch64 | arm64)
    ARCH="arm64"
    ;;
  armv7l)
    ARCH="arm"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get the latest release tag
echo "Fetching the latest release..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not determine the latest release."
    exit 1
fi

echo "Latest release: $LATEST_RELEASE"

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/tax-calculator-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

# Create a temporary directory
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/$EXECUTABLE"

echo "Downloading $DOWNLOAD_URL..."
# Check if the URL exists before downloading
HTTP_CODE=$(curl -sL -w "%{http_code}" -o /dev/null "$DOWNLOAD_URL")
if [ "$HTTP_CODE" != "200" ]; then
    echo "Error: Binary not found at $DOWNLOAD_URL (HTTP $HTTP_CODE)"
    echo "The release may not have been built correctly."
    exit 1
fi

curl -sL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

# Check if installation directory is writable or needs sudo
if [ -w "$INSTALL_DIR" ]; then
    echo "Installing to $INSTALL_DIR/$EXECUTABLE..."
    mv "$TMP_FILE" "$INSTALL_DIR/$EXECUTABLE"
else
    echo "Installing to $INSTALL_DIR/$EXECUTABLE (requires sudo)..."
    sudo mv "$TMP_FILE" "$INSTALL_DIR/$EXECUTABLE"
fi

# Clean up temporary directory
rm -rf "$TMP_DIR"

echo "Installation complete!"
echo "Run '$EXECUTABLE' to start the application."