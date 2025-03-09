#!/bin/bash
set -e

SUDO=${SUDO:-sudo}
if [ "$EUID" -eq 0 ]; then
    SUDO=""
fi

APP_NAME="basement-organizer"
OPT_DEST="/opt/$APP_NAME"
SYMLINK="/usr/local/bin/$APP_NAME"
BUILD_DIR="$HOME/${APP_NAME}_build"

$SUDO apt update && $SUDO apt install -y vim curl git wget unzip ca-certificates

GO_VERSION="go1.22.1"
GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
USR_LOCAL="/usr/local"

if ! command -v go &> /dev/null || [[ "$(go version 2>/dev/null)" != *"${GO_VERSION}"* ]]; then
    echo "Installing Go ${GO_VERSION}..."
    wget -q ${GO_URL} -O go.tar.gz
    $SUDO rm -rf ${USR_LOCAL}/go && $SUDO tar -C ${USR_LOCAL} -xzf go.tar.gz
    rm go.tar.gz
    export PATH=${USR_LOCAL}/go/bin:$PATH
    echo "export PATH=${USR_LOCAL}/go/bin:\$PATH" >> ~/.bashrc
    echo "export PATH=${USR_LOCAL}/go/bin:\$PATH" >> ~/.profile
    source ~/.bashrc
    echo "Go ${GO_VERSION} installed successfully."
else
    echo "Go ${GO_VERSION} is already installed."
fi

if ! command -v google-chrome &> /dev/null; then
    echo "Chrome not found. Installing..."
    wget -q -O google-chrome.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
    $SUDO apt install ./google-chrome.deb -y
    rm google-chrome.deb
else
    echo "Chrome is already installed."
fi

rm -rf "$BUILD_DIR"
echo "Cloning repository into $BUILD_DIR..."
git clone --branch debian-script --single-branch https://github.com/dev-Toumeh/basement-organizer.git "$BUILD_DIR"
cd "$BUILD_DIR" || exit

echo "Installing Go dependencies..."
go mod tidy

if [ ! -f internal/static/js/htmx.min.js ]; then
    echo "Downloading htmx.min.js..."
    mkdir -p internal/static/js
    wget -O internal/static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
fi

echo "Building the Go application..."
go build -o "$APP_NAME" .

echo "Removing existing $OPT_DEST if it exists..."
$SUDO rm -rf "$OPT_DEST"
echo "Moving the repository to $OPT_DEST..."
$SUDO mv "$BUILD_DIR" "$OPT_DEST"

if [ -L "$SYMLINK" ]; then
    $SUDO rm "$SYMLINK"
fi
echo "Creating symbolic link at $SYMLINK..."
$SUDO ln -s "$OPT_DEST/$APP_NAME" "$SYMLINK"

echo "Installation complete. You can now run the application using: $APP_NAME"

