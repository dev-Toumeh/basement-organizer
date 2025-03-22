
# basement-organizer
- [basement-organizer](#basement-organizer)
  - [Overview](#overview)
  - [Installation](#Installation)
  - [Contribution](#contribution)

Overview, detailed information and decisions can be found in [ Technical Documentation ](docs/Technical_Documentation.md).

## Overview
  basement-organizer is an open-source web app designed for local,
   private use to organize your basement using boxes, shelves, and designated areas,
   allows you to search for items without opening every box.
   The app provides a simple, effective solution for keeping your basement orderly and accessible,
   unlike many complex tracking systems built for commercial use.

## Automated Installation
### Linux Debian
 **Step 1**: Install wget (if not already installed)
```bash
  sudo apt update && sudo apt install wget
```
 **Step 2**: Download the installation script
```bash
wget -O install-debian.sh "https://raw.githubusercontent.com/dev-Toumeh/basement-organizer/dev/scripts/install-debian.sh"
```
 **Step 3**: Execute the installation script using sudo
```bash
sudo bash install-debian.sh
```

### Docker
 **Step 1**: Install wget (if not already installed)
```bash
  sudo apt update && sudo apt install wget
```
 **Step 2**: Download the installation scripts
```bash
wget -O install-debian.sh "https://raw.githubusercontent.com/dev-Toumeh/basement-organizer/dev/scripts/install-debian.sh"
wget -O install-docker.sh "https://raw.githubusercontent.com/dev-Toumeh/basement-organizer/dev/scripts/install-docker.sh"
```
 **Step 3**: Execute the installation script using sudo
```bash
sudo bash install-docker.sh
```
## Manual Installation
### Dependencies

Before installing the application, ensure your system has the following dependencies:

Go 1.22.1 or later: The application is built using Go. This dependency is required to compile and execute the program.

Git: The source code is hosted in a Git repository. Git is required to clone the repository and fetch updates.

HTMX: Lightweight JavaScript library that enables dynamic, server-driven user interfaces using HTML attributes.

### Steps


- clone the repository
```bash
git clone https://github.com/dev-Toumeh/basement-organizer.git
```
- download the HTMX library package
   - run the following command in the repository root directory:
    ```bash
        wget -O internal/static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
    ```
   - or download it manually from https://htmx.org/docs/#installing and place it in /path/to/directory/internal/static/js/

- to build the Binary run the following command in the repository root directory:
```bash
go build -o basement -tags prod .
```
- run the application:
```bash
./basement
```
- Launch your browser and navigate to the following URL:
http://localhost:8101

## contribution
### For development see [Getting Started](docs/getting_started.md)
