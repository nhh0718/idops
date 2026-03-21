#!/bin/sh
# idops installer for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/nhh0718/idops/main/install.sh | sh

set -e

REPO="nhh0718/idops"
BINARY="idops"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest version
echo "Fetching latest release..."
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
echo "  Version: $VERSION"

# Download
ASSET="${BINARY}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
TMP_DIR=$(mktemp -d)

echo "Downloading $ASSET..."
curl -fsSL "$URL" -o "$TMP_DIR/$ASSET"

# Extract
echo "Extracting..."
tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
else
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
fi
chmod +x "$INSTALL_DIR/$BINARY"

# Cleanup
rm -rf "$TMP_DIR"

echo ""
echo "idops $VERSION installed to $INSTALL_DIR/$BINARY"
echo "Run: idops --help"
