#!/usr/bin/env bash

# Define variables
APP_NAME="WinGoDarkTray"

# Clean previous builds
echo "Cleaning up..."
if [ -f "$APP_NAME" ]; then
  rm "$APP_NAME"
fi

# Set architecture
GOOS="windows"
GOARCH="amd64"

# Build the Go application
echo "Building the application..."
GOOS=$GOOS GOARCH=$GOARCH go build -ldflags=-H=windowsgui -o "$APP_NAME".exe

echo "Build complete."
