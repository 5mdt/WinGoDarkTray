#!/usr/bin/env bash

# Define variables
APP_NAME="WinGoDarkTray"
ICON_FILE="icon.ico"
BUILD_DIR="./build"

# Function to check if rsrc is installed
check_rsrc_installed() {
  command -v rsrc >/dev/null 2>&1 || {
    echo "rsrc not found, installing...";
    go install github.com/akavel/rsrc@latest;
  }
}

# Function to clean up old build artifacts
clean_up() {
  echo "Cleaning up old artifacts..."
  find . -name "*.exe" -delete
  find $BUILD_DIR -name "*.exe" -delete
  find . -name "*.syso" -delete
}

# Function to create build directory if it doesn't exist
create_build_dir() {
  echo "Ensuring the build directory exists..."
  mkdir -p $BUILD_DIR
}

# Function to embed the icon into the Go application
embed_icon() {
  echo "Embedding the icon..."
  rsrc -ico "$ICON_FILE"
}

# Function to build the Go application for a specific architecture
build_app() {
  local arch=$1
  local arch_name=$2
  echo "Building the application for $arch_name ($arch)..."
  GOARCH=$arch GOOS="windows" go build -ldflags=-H=windowsgui -o "$BUILD_DIR/$APP_NAME-$arch_name.exe"
}

# Function to clean up resources
remove_syso_files() {
  echo "Removing .syso files..."
  find . -name "*.syso" -delete
}

# Main build function
build() {
  clean_up
  create_build_dir
  check_rsrc_installed
  embed_icon
  build_app "amd64" "x64"
  build_app "386" "x32"
  remove_syso_files
  echo "Build complete."
}

# Run the build process
build
