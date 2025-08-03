#!/bin/bash

set -e

REPO="theAnuragMishra/cfgen"
TOOL_NAME="cfgen"
VERSION=${1:-"latest"}

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Normalize OS
case "$OS" in
    Linux*)     GOOS=linux ;;
    Darwin*)    GOOS=darwin ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT) GOOS=windows ;;
    *)          echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Normalize ARCH
case "$ARCH" in
    x86_64|amd64)   GOARCH=amd64 ;;
    arm64|aarch64)  GOARCH=arm64 ;;
    *)              echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Build binary name
EXT=""
[ "$GOOS" = "windows" ] && EXT=".exe"

BINARY_NAME="$TOOL_NAME-$GOOS-$GOARCH$EXT"
TARGET_NAME="$TOOL_NAME$EXT"

# Resolve download URL
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"

echo "Downloading $BINARY_NAME ($VERSION) from $URL..."

curl -L -o "$TARGET_NAME" "$URL"
chmod +x "$TARGET_NAME"

# Optionally move to /usr/local/bin
if [ "$GOOS" != "windows" ]; then
    if [ -w "/usr/local/bin" ]; then
        mv "$TARGET_NAME" /usr/local/bin/$TOOL_NAME
        echo "Installed $TOOL_NAME to /usr/local/bin"
    else
        echo "Move the binary manually:"
        echo "  sudo mv $TARGET_NAME /usr/local/bin/$TOOL_NAME"
    fi
else
    echo "Downloaded $TARGET_NAME. Move it to a folder in your PATH to use it globally."
fi
