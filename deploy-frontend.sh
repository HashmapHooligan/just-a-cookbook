#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$SCRIPT_DIR/frontend/dist/spa"
TARGET_DIR="$HOME/www/kochbuch"

cd "$SCRIPT_DIR/frontend"
echo "Building frontend..."
npm run build

echo "Copying to $TARGET_DIR..."
mkdir -p "$TARGET_DIR"
rm -rf "${TARGET_DIR:?}"/*
cp -r "$DIST_DIR/." "$TARGET_DIR/"

echo "Done. Files at $TARGET_DIR"
