#!/bin/bash

set -e

REPO="smokeyshawn18/up-cli"
VERSION="v1.0.0"
BINARY_NAME="up-cli"

echo "📦 Installing $BINARY_NAME from $REPO (version $VERSION)..."

# Detect OS and ARCH
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  "Linux") PLATFORM="linux" ;;
  "Darwin") PLATFORM="darwin" ;;
  *)
    echo "❌ Unsupported OS: $OS"
    echo "Please download manually from: https://github.com/$REPO/releases/tag/$VERSION"
    exit 1
    ;;
esac

case "$ARCH" in
  "x86_64") ARCH="amd64" ;;
  "arm64" | "aarch64") ARCH="arm64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    echo "Please download manually from: https://github.com/$REPO/releases/tag/$VERSION"
    exit 1
    ;;
esac

ZIP_NAME="${BINARY_NAME}-${PLATFORM}-${ARCH}.zip"
BIN_NAME="${BINARY_NAME}-${PLATFORM}-${ARCH}"

echo "➡️ Downloading $ZIP_NAME..."
curl -LO "https://github.com/$REPO/releases/download/$VERSION/$ZIP_NAME"

echo "📂 Unzipping..."
unzip -o "$ZIP_NAME"

echo "🔧 Making executable..."
chmod +x "$BIN_NAME"

echo "🚚 Moving to /usr/local/bin/$BINARY_NAME"
sudo mv "$BIN_NAME" /usr/local/bin/$BINARY_NAME

echo "✅ Installed! Try running: $BINARY_NAME --help"
