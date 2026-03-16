#!/usr/bin/env bash

set -euo pipefail

REPO="theAnuragMishra/cfgen"
TOOL_NAME="cfgen"
VERSION=${1:-"latest"}

# Detect OS and ARCH
OS_RAW="$(uname -s)"
ARCH_RAW="$(uname -m)"

# Normalize OS
case "$OS_RAW" in
    Linux*)     GOOS=linux ;;
    Darwin*)    GOOS=darwin ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT*) GOOS=windows ;;
    *)          echo "Unsupported OS: $OS_RAW"; exit 1 ;;
esac

# Normalize ARCH
case "$ARCH_RAW" in
    x86_64|amd64)   GOARCH=amd64 ;;
    arm64|aarch64)  GOARCH=arm64 ;;
    *)              echo "Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

# Build binary name
EXT=""
[ "$GOOS" = "windows" ] && EXT=".exe"

BINARY_NAME="$TOOL_NAME-$GOOS-$GOARCH$EXT"
TARGET_NAME="$TOOL_NAME$EXT"

# helper: prefer curl, fallback to wget
fetch() {
    url="$1"
    out="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$out" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$out" "$url"
    else
        echo "Neither curl nor wget is available. Please install one and re-run this script." >&2
        return 2
    fi
}

# Resolve latest tag if requested
if [ "$VERSION" = "latest" ]; then
    echo "Resolving latest release tag..."
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)
    else
        echo "Cannot resolve latest release (curl/wget missing). Please pass a release tag as the first argument." >&2
        exit 2
    fi
    if [ -z "$VERSION" ]; then
        echo "Failed to determine latest version from GitHub API." >&2
        exit 1
    fi
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"

echo "Downloading $BINARY_NAME ($VERSION) from $URL..."

# use a temporary file to avoid partial binary on failure
TMPFILE="$(mktemp -t ${TOOL_NAME}.XXXX 2>/dev/null || mktemp)"
trap 'rm -f "${TMPFILE}"' EXIT

if ! fetch "$URL" "$TMPFILE"; then
    echo "Download failed: $URL" >&2
    exit 1
fi

mv "$TMPFILE" "$TARGET_NAME"
chmod +x "$TARGET_NAME" || true

# Installation target per-platform
if [ "$GOOS" != "windows" ]; then
    DEST_DIR="/usr/local/bin"
    DEST_PATH="$DEST_DIR/$TOOL_NAME"
    if [ -w "$DEST_DIR" ]; then
        mv "$TARGET_NAME" "$DEST_PATH"
        echo "Installed $TOOL_NAME to $DEST_PATH"
    else
        if command -v sudo >/dev/null 2>&1; then
            echo "Installing to $DEST_PATH using sudo (you may be prompted for your password)..."
            if sudo mv "$TARGET_NAME" "$DEST_PATH"; then
                echo "Installed $TOOL_NAME to $DEST_PATH"
            else
                echo "sudo mv failed. You can move the binary manually:" >&2
                echo "  sudo mv $TARGET_NAME $DEST_PATH" 
            fi
        else
            echo "No permission to write to $DEST_DIR and sudo is not available. Move the binary manually:" >&2
            echo "  mv $TARGET_NAME $DEST_PATH"
        fi
    fi
else
    # On Windows, put it under $HOME/bin and suggest adding to PATH
    HOME_BIN="$HOME/bin"
    mkdir -p "$HOME_BIN"
    mv "$TARGET_NAME" "$HOME_BIN/$TARGET_NAME"
    echo "Downloaded $TARGET_NAME to $HOME_BIN/$TARGET_NAME"
    echo "Add $HOME_BIN to your PATH if it's not already in it. For PowerShell add to your profile, or in Git Bash add:"
    echo "  export PATH=\"\$HOME/bin:\$PATH\""
fi

trap - EXIT
echo "Done."
