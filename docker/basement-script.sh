#!/bin/bash
set -e

APP_NAME="basement-organizer"
BIN_DEST="/opt/$APP_NAME"
SYMLINK="/usr/local/bin/$APP_NAME"

# Install required packages
apt update && apt install -y vim curl git wget unzip ca-certificates

GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
GO_VERSION="go1.22.1"
USR_LOCAL="/usr/local"

# Check if Go is installed and matches the expected version
if ! command -v go &> /dev/null || [[ "$(go version 2>/dev/null)" != *"${GO_VERSION}"* ]]; then
    echo "Installing Go ${GO_VERSION}..."
    wget -q ${GO_URL} -O go.tar.gz
    rm -rf ${USR_LOCAL}/go && tar -C ${USR_LOCAL} -xzf go.tar.gz
    rm go.tar.gz
    export PATH=${USR_LOCAL}/go/bin:$PATH
    echo "export PATH=${USR_LOCAL}/go/bin:\$PATH" >> ~/.bashrc
    echo "export PATH=${USR_LOCAL}/go/bin:\$PATH" >> ~/.profile
    source ~/.bashrc
    echo "Go ${GO_VERSION} installed successfully."
else
    echo "Go ${GO_VERSION} is already installed."
fi

# Install Chrome if not installed
if ! command -v google-chrome &> /dev/null; then
    echo "Chrome not found. Installing..."
    wget -q -O google-chrome.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
    apt install ./google-chrome.deb -y
    rm google-chrome.deb
else
    echo "Chrome is already installed."
fi

# Check if the repository exists in /opt
if [ -d "$BIN_DEST" ]; then
    echo "Using existing repository at $BIN_DEST."
else
    echo "Repository not found at $BIN_DEST. Cloning repository..."
    mkdir -p "$BIN_DEST"
    git clone --branch debian-script --single-branch https://github.com/dev-Toumeh/basement-organizer.git "$BIN_DEST"
fi

cd "$BIN_DEST" || exit

echo "Installing Go dependencies..."
go mod tidy

# Ensure internal/static/js exists and copy htmx.min.js if missing
if [ ! -f internal/static/js/htmx.min.js ]; then
    echo "Downloading htmx.min.js..."
    mkdir -p internal/static/js
    wget -O internal/static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
else
    echo "htmx.min.js already exists. Skipping download."
fi

# Check if the binary already exists
if [ -f "$BIN_DEST/$APP_NAME" ]; then
    echo "Deleting the old Binary $BIN_DEST/$APP_NAME."
fi

echo "Building the Go application..."
go build -o "$BIN_DEST/$APP_NAME" .

# Ensure the symbolic link is created
if [ ! -L "$SYMLINK" ]; then
    echo "Creating symbolic link at $SYMLINK..."
    ln -s "$BIN_DEST/$APP_NAME" "$SYMLINK"
else
    echo "Symbolic link $SYMLINK already exists."
fi
echo "Installation complete. You can now run the application using: $APP_NAME"

