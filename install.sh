#!/usr/bin/env bash
set -euo pipefail

TARGET_DIR="${HOME}/.local/bin"
mkdir -p "${TARGET_DIR}"

SCRIPT_DIR=""
if [ -n "${BASH_SOURCE[0]:-}" ] && [ -f "${BASH_SOURCE[0]}" ]; then
	SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi

if [ -n "${SCRIPT_DIR}" ] && [ -f "${SCRIPT_DIR}/main.go" ]; then
	echo "Building palette from local source..."
	go build -ldflags="-s -w" -o "${TARGET_DIR}/palette" "${SCRIPT_DIR}/main.go"
else
	echo "Installing palette from GitHub..."
	TEMP_DIR="$(mktemp -d)"
	trap 'rm -rf "${TEMP_DIR}"' EXIT
	git clone --depth 1 https://github.com/saravenpi/palette.git "${TEMP_DIR}" >/dev/null 2>&1
	(cd "${TEMP_DIR}" && go build -ldflags="-s -w" -o "${TARGET_DIR}/palette" main.go)
fi

if command -v go >/dev/null 2>&1; then
	GOBIN_DIR="$(go env GOPATH)/bin"
	if [ -d "${GOBIN_DIR}" ] && [ "${GOBIN_DIR}" != "${TARGET_DIR}" ]; then
		cp "${TARGET_DIR}/palette" "${GOBIN_DIR}/palette"
	fi
fi

case ":${PATH}:" in
	*:"${TARGET_DIR}":*) ;;
	*) echo "Note: Add ${TARGET_DIR} to your PATH to run 'palette' directly." ;;
esac

echo "Installed palette to ${TARGET_DIR}/palette"
"${TARGET_DIR}/palette" --version
