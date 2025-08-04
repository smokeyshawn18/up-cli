#!/bin/bash

set -e

APP_NAME="up-cli"
DIST="dist"
ENTRY="./cmd/app/main.go"
VERSION="v1.0.0"

echo "🔨 Building $APP_NAME version $VERSION for multiple platforms..."

# Create dist directory
mkdir -p $DIST

# Clean previous builds
rm -f $DIST/*

# Build matrix
build() {
  GOOS=$1
  GOARCH=$2
  EXT=$3
  OUTPUT="${DIST}/${APP_NAME}-${GOOS}-${GOARCH}${EXT}"

  echo "📦 Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=${VERSION}" -o "$OUTPUT" "$ENTRY"
}

# Linux builds
build linux amd64 ""
build linux arm64 ""

# macOS builds
build darwin amd64 ""
build darwin arm64 ""

# Windows builds
build windows amd64 ".exe"
build windows arm64 ".exe"

# Zip all binaries
echo "🗜️ Zipping binaries..."
cd $DIST
for f in *; do
  zip -q "${f}.zip" "$f"
done
cd ..

echo "✅ Build complete. Binaries and archives are in ./$DIST"
