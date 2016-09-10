#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Fetches Raspbian Jessie Lite and flash it to an SDCard.
# Then it updates the SDCard so it automatically self-initializes as a dlibox on
# first boot.

set -eu
cd "$(dirname $0)"

if [ "$#" -ne 2 ]; then
  echo "Flashes a customized version of Raspbian that automatically"
  echo "sets itself up as a dlibox."
  echo ""
  echo "usage: ./flash.sh /dev/<sdcard_path> <ssid>"
  exit 1
fi

SDCARD=$1
SSID="$2"

echo "Warning! This will blow up everythingin in $SDCARD"
echo ""
echo "This script has minimal use of 'sudo' for 'dd' and modifying the partitions"
echo ""

./prep/flash.sh $SDCARD
./prep/edit.sh $SDCARD "$SSID"

echo ""
echo "You can now remove the SDCard safely and boot your Raspberry Pi"
echo "Once it was booted for several seconds, you can find the hostname with:"
echo "  ./find.sh"
echo ""
echo "Then connect with:"
echo "  ssh -o StrictHostKeyChecking=no pi@raspberrypi"
echo ""
echo "You can follow the update process by either connecting a monitor"
echo "to the HDMI port or by ssh'ing into the Pi and running:"
echo "  tail -f /var/log/dlibox_firstboot.log"
