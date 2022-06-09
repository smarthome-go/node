#!/bin/bash

# Checks
echo -e "\x1b[1;34mInstallation running...\x1b[1;0m"
echo "Checking system architecture..."
grep -q BCM /proc/cpuinfo || ( echo -e "\x1b[1;31mPlease run the installer on a Raspberry-Pi\x1b[1;0m" && exit 1 )
echo "Detected system architecture: BCM / Raspberry-Pi"

# Installation
echo "Creating application root..."
sudo mkdir -p /usr/bin/smarthome-hw || exit 1
echo "Changing ownership of application root..."
sudo chown -R pi /usr/bin/smarthome-hw  || exit 1
echo "Moving application files to new application root..."
mv ./smarthome-hw /usr/bin/smarthome-hw/ || exit 1
echo "Installing systemd service..."
sudo mv ./smarthome-hw.service /lib/systemd/system/smarthome-hw.service || exit 1

# Reload systemd
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload || exit 1
echo "Activating service using systemd..."
sudo systemctl start smarthome-hw || exit 1

# Finishing
echo -e "\x1b[1;32mInstallation completed\x1b[1;0m"
