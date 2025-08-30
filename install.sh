#!/usr/bin/env bash
set -euo pipefail

REPO="unholyFigaro/shh"
APP="shh"
VERSION="${SHH_VERSION:-latest}"

# --- detect OS ---
uname_s="$(uname -s)"
case "$uname_s" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *) echo "❌ Unsupported OS: $uname_s"; exit 1 ;;
esac

# --- detect ARCH ---
uname_m="$(uname -m)"
case "$uname_m" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  armv7l|armv7) ARCH="armv7" ;;
  *) echo "❌ Unsupported ARCH: $uname_m"; exit 1 ;;
esac

# --- choose install dir ---
DEFAULT_DIR="/usr/local/bin"
TARGET_DIR="${SHH_INSTALL_DIR:-$DEFAULT_DIR}"

# if default dir not writable, fall back to ~/.local/bin
if [ ! -w "$TARGET_DIR" ]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    TARGET_DIR="$HOME/.local/bin"
    mkdir -p "$TARGET_DIR"
    SUDO=""
  fi
else
  SUDO=""
fi

# --- build download URLs ---
if [ "$VERSION" = "latest" ]; then
  BASE="https://github.com/$REPO/releases/latest/download"
else
  # ensure leading 'v' is accepted either way
  case "$VERSION" in v*) TAG="$VERSION" ;; *) TAG="v$VERSION" ;; esac
  BASE="https://github.com/$REPO/releases/download/$TAG"
fi

ASSET="${APP}_${OS}_${ARCH}.tar.gz"
URL="$BASE/$ASSET"

# --- temp dir ---
TMPDIR="$(mktemp -d)"
cleanup() { rm -rf "$TMPDIR"; }
trap cleanup EXIT

echo "⬇️  Downloading $APP ($OS/$ARCH) from: $URL"
if ! curl -fsSL -o "$TMPDIR/$ASSET" "$URL"; then
  echo "⚠️  Prebuilt binary not found or download failed."
  if command -v go >/dev/null 2>&1; then
    echo "🔁 Falling back to 'go install github.com/$REPO@latest'..."
    GO_BIN_DIR="$(go env GOPATH 2>/dev/null)/bin"
    GO_BIN_DIR="${GO_BIN_DIR:-$HOME/go/bin}"
    GOMODCACHE="$(go env GOMODCACHE 2>/dev/null || true)"
    # Try installing
    go install "github.com/$REPO@latest"
    # Copy to TARGET_DIR if needed
    if [ -x "$GO_BIN_DIR/$APP" ]; then
      if [ -n "$SUDO" ]; then
        $SUDO install -m 0755 "$GO_BIN_DIR/$APP" "$TARGET_DIR/$APP"
      else
        install -m 0755 "$GO_BIN_DIR/$APP" "$TARGET_DIR/$APP"
      fi
      echo "✅ Installed $APP to $TARGET_DIR/$APP"
      if ! command -v "$APP" >/dev/null 2>&1; then
        echo "ℹ️  Add to PATH: export PATH=\"$TARGET_DIR:\$PATH\""
      fi
      exit 0
    else
      echo "❌ go install finished but $GO_BIN_DIR/$APP not found."
      exit 1
    fi
  else
    echo "❌ Go is not installed, and prebuilt binaries are unavailable."
    echo "   Install Go or build from source manually."
    exit 1
  fi
fi

echo "📦 Extracting archive..."
tar -xzf "$TMPDIR/$ASSET" -C "$TMPDIR"

# locate the binary inside the tarball
BIN_PATH="$(find "$TMPDIR" -type f -name "$APP" -perm -u+x | head -n1 || true)"
if [ -z "$BIN_PATH" ]; then
  # fallback: sometimes it's in 'dist' or similar
  BIN_PATH="$(find "$TMPDIR" -type f -name "$APP" | head -n1 || true)"
fi
if [ -z "$BIN_PATH" ]; then
  echo "❌ Could not find '$APP' binary inside the archive."
  exit 1
fi

echo "🛠  Installing to $TARGET_DIR (may require sudo)..."
if [ -n "$SUDO" ]; then
  $SUDO install -m 0755 "$BIN_PATH" "$TARGET_DIR/$APP"
else
  install -m 0755 "$BIN_PATH" "$TARGET_DIR/$APP"
fi

echo "✅ Installed $APP to $TARGET_DIR/$APP"

if ! command -v "$APP" >/dev/null 2>&1; then
  # If the installed dir is not in PATH, suggest adding it
  case ":$PATH:" in
    *":$TARGET_DIR:"*) ;;
    *) echo "ℹ️  Add to PATH: echo 'export PATH=\"$TARGET_DIR:\$PATH\"' >> ~/.bashrc && source ~/.bashrc" ;;
  esac
fi

echo
echo "🚀 Try it:"
echo "   $APP --help"
echo "   $APP add dev-stand --host 127.0.0.1 -p 2222 -u $(id -un)"