#!/usr/bin/env bash
# DPTR Install Script
# Builds the Go binary and sets up the systemd user service.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY_NAME="dptr"
INSTALL_BIN="$HOME/.local/bin"
INSTALL_CONFIG="$HOME/.config/dptr"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"

echo "==> Building DPTR..."
cd "$PROJECT_DIR"

# Prefer nix-shell if go isn't in PATH (NixOS users)
if ! command -v go &>/dev/null; then
  echo "    'go' not in PATH; using nix-shell..."
  nix-shell -p go --run "go build -o $BINARY_NAME ./cmd/dptr"
else
  go build -o "$BINARY_NAME" ./cmd/dptr
fi

echo "==> Installing binary to $INSTALL_BIN/$BINARY_NAME"
mkdir -p "$INSTALL_BIN"
cp "$BINARY_NAME" "$INSTALL_BIN/$BINARY_NAME"
chmod +x "$INSTALL_BIN/$BINARY_NAME"

echo "==> Installing config template to $INSTALL_CONFIG/"
mkdir -p "$INSTALL_CONFIG"
if [ ! -f "$INSTALL_CONFIG/config.yaml" ]; then
  cp "$PROJECT_DIR/config.template.yaml" "$INSTALL_CONFIG/config.yaml"
  echo "    Copied config.template.yaml -> $INSTALL_CONFIG/config.yaml"
  echo "    !! Please edit $INSTALL_CONFIG/config.yaml before first run !!"
else
  echo "    Config already exists at $INSTALL_CONFIG/config.yaml — skipping."
fi

echo "==> Installing systemd user service..."
mkdir -p "$SYSTEMD_USER_DIR"
# Patch the service file to point at the actual project dir for modules/
SERVICE_SRC="$SCRIPT_DIR/dptr.service"
SERVICE_DEST="$SYSTEMD_USER_DIR/dptr.service"
cp "$SERVICE_SRC" "$SERVICE_DEST"

systemctl --user daemon-reload
systemctl --user enable dptr.service
echo "    Service enabled."

echo ""
echo "╔══════════════════════════════════════════════════════╗"
echo "║  DPTR installed successfully!                        ║"
echo "║                                                      ║"
echo "║  Next steps:                                         ║"
echo "║  1. Edit ~/.config/dptr/config.yaml                  ║"
echo "║  2. Add API keys to ~/.config/dptr/.env              ║"
echo "║  3. Test: dptr --force --config ~/.config/dptr/config.yaml ║"
echo "║  4. The report will auto-appear on next login         ║"
echo "╚══════════════════════════════════════════════════════╝"
