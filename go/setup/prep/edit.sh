#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Edits a Raspbian Jessie Lite disk so it automatically self-initializes as a
# dlibox on first boot.

set -eu
cd "$(dirname $0)"


if [ "$#" -ne 2 ]; then
  echo "Customize Raspbian to automatically sets itself up as a dlibox."
  echo ""
  echo "usage: ./edit.sh /dev/<sdcard_path> <ssid>"
  exit 1
fi


# TODO(maruel): Some confirmation or verification. A user could destroy their
# workstation easily.
# Linux generally use /dev/sdX, OSX uses /dev/diskN.
SDCARD=$1
SSID="$2"


echo "- Mounting"
./umount.sh $SDCARD  &>/dev/null
# Needs 'p' for /dev/mmcblkN but not for /dev/sdX
BOOT=$(LANG=C /usr/bin/udisksctl mount -b ${SDCARD}*1 | sed 's/.\+ at \(.\+\)\+\./\1/')
echo "- /boot mounted as $BOOT"
ROOT=$(LANG=C /usr/bin/udisksctl mount -b ${SDCARD}*2 | sed 's/.\+ at \(.\+\)\+\./\1/')
echo "- / mounted as $ROOT"


# Skip this if you don't use a small display.
# Strictly speaking, you won't need a monitor at all since ssh will be up and
# running and the device will connect to the SSID provided.
if [ false ]; then
  echo "- Enabling 5\" display support (optional)"
  ./enable_5inch.sh $BOOT
fi


# Setup SSH keys, wifi and automatic setup process on first boot.
if [ true ]; then
  ./setup_firstboot.sh $ROOT "$SSID"
fi


echo "- Unmounting"
sync
./umount.sh $SDCARD
