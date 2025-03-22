#!/bin/bash
set -e

APP_NAME="basement-organizer"
LOCAL_SHARE="$HOME/.local/share"
REPO_DIR="$LOCAL_SHARE/$APP_NAME"
SYMLINK="/usr/local/bin/$APP_NAME"

# Ensure directory exists
mkdir -p "$LOCAL_SHARE"

# Install dependencies
sudo apt update && sudo apt install -y curl git wget unzip

# Install Go if not present or wrong version
GO_VERSION="go1.22.1"
GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
if ! command -v go &> /dev/null || [[ "$(go version 2>/dev/null)" != *"${GO_VERSION}"* ]]; then
    echo "Installing Go ${GO_VERSION}..."
    wget -q ${GO_URL} -O go.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go.tar.gz
    rm go.tar.gz
    export PATH=/usr/local/go/bin:$PATH
    echo "export PATH=/usr/local/go/bin:\$PATH" >> ~/.bashrc
    echo "export PATH=/usr/local/go/bin:\$PATH" >> ~/.profile
    source ~/.bashrc
fi

# Clone or update repo
if [ -d "$REPO_DIR/.git" ]; then
    echo "Repository already exists. Pulling latest changes..."
    git -C "$REPO_DIR" pull
else
    echo "Cloning repository..."
    git clone https://github.com/dev-Toumeh/basement-organizer.git "$REPO_DIR"
fi

cd "$REPO_DIR"

echo "Tidying Go modules..."
go mod tidy

if [ ! -f internal/static/js/htmx.min.js ]; then
    echo "Installing htmx.min.js..."
    mkdir -p internal/static/js
    wget -O internal/static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
fi

echo "Building the application..."
go build -o "$APP_NAME" -tags prod .

# Create symlink
if [ -L "$SYMLINK" ]; then
    sudo rm "$SYMLINK"
fi
sudo ln -s "$REPO_DIR/$APP_NAME" "$SYMLINK"

echo "Done. You can now run the app using: $APP_NAME"
