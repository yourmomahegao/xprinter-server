#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "ERROR: Run as root required"
  exit 1
fi

echo "Creating directory..."
mkdir -p /opt/xprinter
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to create /opt/xprinter"
  exit 1
fi

echo "Copying binary..."
cp xprinter /opt/xprinter/
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to copy xprinter"
  exit 1
fi

chmod +x /opt/xprinter/xprinter
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to set executable permissions"
  exit 1
fi

echo "Installing service..."
cp xprinter.service /etc/systemd/system/
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to copy service file"
  exit 1
fi

echo "Reloading systemd..."
systemctl daemon-reload
if [ $? -ne 0 ]; then
  echo "ERROR: systemctl daemon-reload failed"
  exit 1
fi

echo "Enabling service..."
systemctl enable xprinter
if [ $? -ne 0 ]; then
  echo "ERROR: systemctl enable failed"
  exit 1
fi

echo "Starting service..."
systemctl restart xprinter
if [ $? -ne 0 ]; then
  echo "ERROR: systemctl restart failed"
  exit 1
fi

echo "======================================"
echo "SUCCESS: XPrinter installed and running"
echo "Service: xprinter"
echo "Path: /opt/xprinter/xprinter"
echo "======================================"