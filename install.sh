#!/usr/bin/env bash

set -e

REPO="amintehrani/llm-gate"
BINARY_NAME="llm-gate"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_NAME="linux";;
    Darwin*)    OS_NAME="darwin";;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)   ARCH_NAME="amd64";;
    arm64)    ARCH_NAME="arm64";;
    aarch64)  ARCH_NAME="arm64";;
    *)        echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

echo "Detected OS: ${OS_NAME}, Architecture: ${ARCH_NAME}"

# Fetch latest release data
LATEST_RELEASE_URL="https://api.github.com/repos/${REPO}/releases/latest"
echo "Fetching latest release information..."

# Extract the browser download url for the correct OS and Architecture
DOWNLOAD_URL=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url" | grep "${OS_NAME}_${ARCH_NAME}.tar.gz" | cut -d '"' -f 4 | head -n 1 || true)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find a release for ${OS_NAME}_${ARCH_NAME}."
    echo "Please check the releases page: https://github.com/${REPO}/releases"
    exit 1
fi

echo "Downloading ${BINARY_NAME}..."
TMP_DIR=$(mktemp -d)
# Clean up temp directory on exit
trap 'rm -rf "$TMP_DIR"' EXIT

curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/${BINARY_NAME}.tar.gz"

echo "Extracting..."
tar -xzf "$TMP_DIR/${BINARY_NAME}.tar.gz" -C "$TMP_DIR"

# Determine install directory
if [ -d "$HOME/.local/bin" ] && [ -w "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    USE_SUDO=""
else
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ]; then
        echo "Requires sudo privileges to install to $INSTALL_DIR"
        USE_SUDO="sudo"
    else
        USE_SUDO=""
    fi
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
$USE_SUDO mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
$USE_SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo
echo "✅ $BINARY_NAME was successfully installed to $INSTALL_DIR/$BINARY_NAME!"
echo "Run '$BINARY_NAME --help' to get started."

# Check if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo
    echo "⚠️  WARNING: $INSTALL_DIR is not in your PATH."
    echo "Please add it to your profile (e.g., ~/.bashrc or ~/.zshrc):"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
fi
