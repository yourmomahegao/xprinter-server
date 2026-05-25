#!/bin/bash

set -e

echo "Cleaning..."
rm -rf build
mkdir -p build/windows
mkdir -p build/linux

echo "Building Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o build/windows/xprinter.exe ./cmd/xprinter/xprinter_windows_service

echo "Building Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -o build/windows/xprinter_arm64.exe ./cmd/xprinter/xprinter_windows_service

echo "Building Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o build/linux/xprinter ./cmd/xprinter/xprinter_linux_service

echo "Building Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o build/linux/xprinter_arm64 ./cmd/xprinter/xprinter_linux_service

echo "Copying deploy scripts..."

cp deploy/windows/install.bat build/windows/
cp deploy/linux/install.sh build/linux/
cp deploy/linux/xprinter.service build/linux/

chmod +x build/linux/install.sh

echo "DONE"