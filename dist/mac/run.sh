#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"
ARCH="$(uname -m)"
echo "=== Starting Loadin for macOS ($ARCH) ==="
if [ "$ARCH" = "arm64" ]; then
    chmod +x "$DIR/Loadin_darwin_arm64" 2>/dev/null
    exec "$DIR/Loadin_darwin_arm64" "$@"
else
    chmod +x "$DIR/Loadin_darwin_amd64" 2>/dev/null
    exec "$DIR/Loadin_darwin_amd64" "$@"
fi