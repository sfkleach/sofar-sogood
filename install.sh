#!/bin/bash
# Install the latest So far, so good release from GitHub.
# Usage: install.sh [--to DIR]   (default DIR is /usr/local/bin)

set -euo pipefail

DEST="/usr/local/bin"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --to)
            DEST="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

if [[ ! -d "$DEST" ]]; then
    echo "Error: Destination directory $DEST does not exist." >&2
    exit 1
fi

REPO="sfkleach/sofar-sogood"
LATEST_RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"

ARCH=$(uname -m)
OS=$(uname | tr '[:upper:]' '[:lower:]')

# The platform name is the suffix of the release asset, which has the form
# sofar-sogood_<tag>_<platform>.tar.gz (or .zip on Windows).
BINARY="sofar-sogood"
case "$OS" in
    linux)
        case "$ARCH" in
            x86_64) PLATFORM="linux-amd64.tar.gz" ;;
            aarch64) PLATFORM="linux-arm64.tar.gz" ;;
            *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
        esac
        ;;
    darwin)
        case "$ARCH" in
            x86_64) PLATFORM="macos-amd64.tar.gz" ;;
            arm64) PLATFORM="macos-arm64.tar.gz" ;;
            *) echo "Unsupported Mac architecture: $ARCH" >&2; exit 1 ;;
        esac
        ;;
    mingw*|msys*|cygwin*)
        PLATFORM="windows-amd64.zip"
        BINARY="sofar-sogood.exe"
        ;;
    *)
        echo "Unsupported OS: $OS" >&2
        exit 1
        ;;
esac

TMP_DIR=$(mktemp -d)
# The path is a fresh mktemp directory created above, so removing it is safe.
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Fetching latest So far, so good release..."
ASSET_URL=$(curl -fsS "$LATEST_RELEASE_URL" | grep "browser_download_url" | cut -d '"' -f 4 | grep "_$PLATFORM\$" || true)

if [[ -z "$ASSET_URL" ]]; then
    echo "Error: Unable to find a So far, so good release for your system." >&2
    exit 1
fi

echo "Downloading $ASSET_URL..."
curl -fL -o "$TMP_DIR/archive" "$ASSET_URL"

echo "Extracting files..."
if [[ "$PLATFORM" == *.tar.gz ]]; then
    tar -xzf "$TMP_DIR/archive" -C "$TMP_DIR"
else
    unzip -q "$TMP_DIR/archive" -d "$TMP_DIR"
fi

# The archive contains a single directory holding the binary alongside the README and example config.
BIN_PATH=$(find "$TMP_DIR" -name "$BINARY" -type f | head -n 1)
if [[ -z "$BIN_PATH" ]]; then
    echo "Error: $BINARY not found in the downloaded archive." >&2
    exit 1
fi
chmod +x "$BIN_PATH"

if [[ -w "$DEST" ]]; then
    mv "$BIN_PATH" "$DEST/"
else
    echo "No write permissions for $DEST, using sudo..."
    sudo mv "$BIN_PATH" "$DEST/"
fi

echo "Installation complete! $BINARY is installed to $DEST."
if [[ "$OS" == "darwin" ]]; then
    echo "Note: macOS may block the binary because it is not notarized."
    echo "      See the release notes, or run: xattr -rd com.apple.quarantine \"$DEST/$BINARY\""
fi
