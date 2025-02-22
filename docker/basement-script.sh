#!/bin/bash

APP_NAME="basement-organizer"
INSTALL_DIR="/opt"
SYMLINK="/usr/local/bin/$APP_NAME"

# Install required packages
apt update && apt install -y vim curl git wget unzip ca-certificates npm 
npm install -g nodemon

# Install Go if not installed or outdated
GO_VERSION="go1.22.1"
GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
USR_LOCAL="/usr/local"

if ! command -v go &> /dev/null || [[ $(go version) != *"${GO_VERSION}"* ]]; then
    echo "Installing Go ${GO_VERSION}..."
    wget -q ${GO_URL} -O go.tar.gz
    rm -rf ${USR_LOCAL}/go && tar -C ${USR_LOCAL} -xzf go.tar.gz
    rm go.tar.gz
    export PATH=${USR_LOCAL}/go/bin:$PATH
    echo "export PATH=${USR_LOCAL}/go/bin:\$PATH" >> ~/.bashrc
    source ~/.bashrc
else
    echo "Go ${GO_VERSION} is already installed."
fi

# Install Chrome if not installed
if ! command -v google-chrome &> /dev/null; then
    echo "Chrome not found. Installing..."
    wget -q -O google-chrome.deb https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
    dpkg -i google-chrome.deb || apt-get -f install -y
    rm google-chrome.deb
else
    echo "Chrome is already installed."
fi

# Clone the repository if it doesn't exist
if [ ! -d "$APP_NAME" ]; then
    echo "Cloning the repository..."
    git clone --branch debian-script --single-branch https://github.com/dev-Toumeh/basement-organizer.git
    # git clone https://github.com/dev-Toumeh/basement-organizer.git
else
    echo "Repository already exists. Skipping cloning."
fi

# Change into the project directory
cd "$APP_NAME" || exit

# Install Node.js dependencies if package.json exists
if [ -f package.json ]; then
    echo "Installing Node.js dependencies..."
    npm install
fi

# Install Go dependencies if go.mod exists
echo "Installing Go dependencies..."
go mod tidy

# Ensure internal/static/js directory exists
mkdir -p internal/static/js

# Copy htmx.min.js if not already present
if [ ! -f internal/static/js/htmx.min.js ]; then
    echo "Copying htmx.min.js..."
    cp ./node_modules/htmx.org/dist/htmx.min.js internal/static/js/htmx.min.js
else
    echo "htmx.min.js already exists. Skipping copy."
fi

# Build the Go application
if [ -f main.go ]; then
    echo "Building the Go application..."
    go build -o $APP_NAME
fi

# Move the application to /opt
echo "Moving application to $INSTALL_DIR..."
mkdir -p $INSTALL_DIR
cd .. && mv $APP_NAME $INSTALL_DIR/

# Create a symbolic link for easy execution
echo "Creating symbolic link at $SYMLINK..."
ln -sf $INSTALL_DIR/$APP_NAME/$APP_NAME $SYMLINK

echo "Installation complete."
echo "You can now run the application using: $APP_NAME"
