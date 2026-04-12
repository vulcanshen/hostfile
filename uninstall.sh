#!/bin/sh
# hostfile uninstaller for macOS / Linux / Git Bash
# Usage: curl -fsSL https://raw.githubusercontent.com/vulcanshen/hostfile/main/uninstall.sh | sh

set -e

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  mingw*|msys*|cygwin*) OS="windows" ;;
esac

# Determine install locations to check
if [ "$OS" = "windows" ]; then
  CANDIDATES="$HOME/bin/hostfile.exe"
else
  CANDIDATES="$HOME/.local/bin/hostfile /usr/local/bin/hostfile"
fi

FOUND=""
for path in $CANDIDATES; do
  if [ -f "$path" ]; then
    FOUND="$path"
    break
  fi
done

if [ -z "$FOUND" ]; then
  echo "hostfile not found in expected locations."
  echo "Checked: $CANDIDATES"
  exit 1
fi

rm "$FOUND"
echo "removed $FOUND"

# Remove saved snapshots if present
SAVE_DIR="$HOME/.hostfile"
if [ -d "$SAVE_DIR" ]; then
  printf "Remove saved snapshots in %s? [y/N]: " "$SAVE_DIR"
  read -r answer
  case "$answer" in
    y|Y|yes|YES)
      rm -rf "$SAVE_DIR"
      echo "removed $SAVE_DIR"
      ;;
    *)
      echo "kept $SAVE_DIR"
      ;;
  esac
fi

echo ""
echo "hostfile uninstalled."
