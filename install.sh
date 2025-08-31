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
  echo "❌ Release asset not found or download failed."
  echo "   Expected asset: $ASSET  (version: ${VERSION})"
  exit 1
fi

echo "📦 Extracting archive..."
tar -xzf "$TMPDIR/$ASSET" -C "$TMPDIR"

# locate the binary inside the tarball
BIN_PATH="$(find "$TMPDIR" -type f -name "$APP" -perm -u+x | head -n1 || true)"
if [ -z "$BIN_PATH" ]; then
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
  case ":$PATH:" in
    *":$TARGET_DIR:"*) ;;
    *) echo "ℹ️  Add to PATH and restart shell:  export PATH=\"$TARGET_DIR:\$PATH\"" ;;
  esac
fi

# -------------------------
#   Completion installer
# -------------------------
detect_shell() {
  # allow override: SHH_COMPLETION=bash|zsh|fish|all|none
  if [ -n "${SHH_COMPLETION:-}" ]; then
    echo "$SHH_COMPLETION"
    return
  fi
  # fallback to current shell basename
  b="$(basename "${SHELL:-}")"
  case "$b" in
    bash|zsh|fish) echo "$b" ;;
    *) echo "none" ;;
  esac
}

install_completion_bash() {
  local dir="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
  mkdir -p "$dir"
  if ! "$TARGET_DIR/$APP" completion bash > "$dir/$APP" 2>/dev/null; then
    echo "⚠️  Failed to generate bash completion."
    return
  fi
  echo "✅ Bash completion installed → $dir/$APP"
  echo "   Make sure bash-completion is loaded. If not, add to ~/.bashrc:"
  echo "     [[ -f /usr/share/bash-completion/bash_completion ]] && . /usr/share/bash-completion/bash_completion"
}

install_completion_zsh() {
  local dir="${ZSH_COMPLETIONS_DIR:-$HOME/.zsh/completions}"
  mkdir -p "$dir"
  if ! "$TARGET_DIR/$APP" completion zsh > "$dir/_$APP" 2>/dev/null; then
    echo "⚠️  Failed to generate zsh completion."
    return
  fi
  echo "✅ Zsh completion installed → $dir/_$APP"
  echo "   If completion is not active, add to ~/.zshrc and restart zsh:"
  echo "     fpath+=(~/.zsh/completions)"
  echo "     autoload -U compinit && compinit"
}

install_completion_fish() {
  local dir="$HOME/.config/fish/completions"
  mkdir -p "$dir"
  if ! "$TARGET_DIR/$APP" completion fish > "$dir/$APP.fish" 2>/dev/null; then
    echo "⚠️  Failed to generate fish completion."
    return
  fi
  echo "✅ Fish completion installed → $dir/$APP.fish"
}

install_completion() {
  local which="${1:-$(detect_shell)}"
  case "$which" in
    none)
      echo "ℹ️  Skip shell completion (unknown shell). You can install manually:"
      echo "    $APP completion bash|zsh|fish"
      ;;
    all)
      install_completion_bash || true
      install_completion_zsh || true
      install_completion_fish || true
      ;;
    bash) install_completion_bash ;;
    zsh)  install_completion_zsh ;;
    fish) install_completion_fish ;;
    *)
      echo "ℹ️  Unknown SHH_COMPLETION=$which. Supported: bash|zsh|fish|all|none"
      ;;
  esac
}

echo "🔧 Setting up shell completion..."
install_completion

echo
echo "🚀 Try it:"
echo "   $APP --help"
echo "   $APP add dev-stand --host 127.0.0.1 -p 2222 -u \$(id -un)"
echo
echo "💡 Tip: open a NEW terminal to activate completion, or source your rc file."
