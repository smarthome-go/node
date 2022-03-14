#!/bin/bash

grep -q BCM /proc/cpuinfo || ( echo -e "\x1b[1;31mPlease run the installer on a Raspberry pi\x1b[1;0m" && exit 1 )

echo -e "\x1b[1;34minstallation running\x1b[1;0m"

sudo mkdir -p /usr/bin/smarthome-hw || exit 1
sudo chown -R pi /usr/bin/smarthome-hw  || exit 1
mv ./smarthome-hw /usr/bin/smarthome-hw/ || exit 1
sudo cp ./smarthome-hw.service /lib/systemd/system/smarthome-hw.service || exit 1

# Reload systemd
sudo systemctl daemon-reload || exit 1
sudo systemctl start smarthome-hw || exit 1

echo -e "\x1b[1;32minstallation completed\x1b[1;0m"
