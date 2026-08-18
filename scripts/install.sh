#!/usr/bin/env sh
# claude-menu installer for macOS / Linux.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/OWNER/claude-menu/main/scripts/install.sh | sh
# Override the source repo or version:
#   CLAUDE_MENU_REPO=OWNER/claude-menu VERSION=v1.0.0 sh install.sh
set -eu

REPO="${CLAUDE_MENU_REPO:-OWNER/claude-menu}"
VERSION="${VERSION:-latest}"
BINARY="claude-menu"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $os"; exit 1 ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "Unsupported arch: $arch"; exit 1 ;;
esac

asset="${BINARY}_${os}_${arch}"

# Pick a writable install dir.
if [ -w /usr/local/bin ] 2>/dev/null; then
  dir="/usr/local/bin"
else
  dir="$HOME/.local/bin"
fi
mkdir -p "$dir"

if [ "$VERSION" = "latest" ]; then
  url="https://github.com/$REPO/releases/latest/download/$asset"
else
  url="https://github.com/$REPO/releases/download/$VERSION/$asset"
fi

echo "Downloading $url"
tmp="$(mktemp)"
if curl -fsSL "$url" -o "$tmp"; then
  chmod +x "$tmp"
  mv "$tmp" "$dir/$BINARY"
  echo "Installed $dir/$BINARY"
elif command -v go >/dev/null 2>&1 && [ -f go.mod ]; then
  rm -f "$tmp"
  echo "Download failed — building from source..."
  go build -ldflags "-s -w" -o "$dir/$BINARY" .
  echo "Installed $dir/$BINARY"
else
  rm -f "$tmp"
  echo "Download failed and no Go toolchain / source available." >&2
  echo "Set CLAUDE_MENU_REPO / VERSION correctly, or run from the source tree with Go installed." >&2
  exit 1
fi

case ":$PATH:" in
  *":$dir:"*) ;;
  *) echo "NOTE: add $dir to your PATH, e.g. echo 'export PATH=\"$dir:\$PATH\"' >> ~/.zshrc" ;;
esac

echo "Run: $BINARY"
