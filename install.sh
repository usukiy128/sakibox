#!/usr/bin/env sh
set -e

REPO="usukiy128/sakibox"
BIN_NAME="sakibox"
INSTALL_DIR="${SAKIBOX_INSTALL_DIR:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin|Linux) ;;
  *)
    echo "Unsupported OS"
    exit 1
    ;;
esac

if ! command -v git >/dev/null 2>&1; then
  echo "git is required"
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "go is required"
  exit 1
fi

TMP_DIR=$(mktemp -d)
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

git clone --depth 1 "https://github.com/${REPO}.git" "$TMP_DIR/sakibox" >/dev/null 2>&1

cd "$TMP_DIR/sakibox"

go build -o "$BIN_NAME"

if [ ! -d "$INSTALL_DIR" ]; then
  if [ -w "$(dirname "$INSTALL_DIR")" ]; then
    mkdir -p "$INSTALL_DIR"
  else
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
else
  sudo install -m 0755 "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
fi

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
