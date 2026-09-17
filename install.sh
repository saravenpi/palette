#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="${HOME}/.local/bin"

mkdir -p "${TARGET_DIR}"

echo "Building palette..."
go build -ldflags="-s -w" -o "${TARGET_DIR}/palette" "${REPO_DIR}/main.go"

echo "Installed palette to ${TARGET_DIR}/palette"
"${TARGET_DIR}/palette" --version
