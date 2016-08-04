#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Fetches Raspbian Jessie Lite and flashes it to an SDCard.

# TODO(someone): Make this script OSX compatible. For now it was only tested on
# Ubuntu.

set -eu
cd "$(dirname $0)"


if [ "$#" -ne 1 ]; then
  echo "Flashes Raspbian to an SD card"
  echo ""
  echo "usage: ./flash.sh /dev/<sdcard_path>"
  exit 1
fi


# TODO(maruel): Some confirmation or verification. A user could destroy their
# workstation easily.
# Linux generally use /dev/sdX, OSX uses /dev/diskN.
SDCARD=$1


echo "- Unmounting"
./umount.sh $SDCARD &>/dev/null


# TODO(maruel): Figure the name automatically.
IMGNAME=2016-05-27-raspbian-jessie-lite.img
if [ ! -f $IMGNAME ]; then
  echo "- Fetching Raspbian Jessie Lite latest"
  curl -L -o raspbian_lite_latest.zip https://downloads.raspberrypi.org/raspbian_lite_latest
  unzip raspbian_lite_latest.zip
fi


echo "- Flashing (takes 2 minutes)"
sudo /bin/bash -c "time dd bs=4M if=$IMGNAME of=$SDCARD"
echo "- Flushing I/O cache"
# This is important otherwise the mount afterward may 'see' the old partition
# table.
time sync


echo "- Reloading partition table"
# Wait a bit to try to workaround "Error looking up object for device" when
# immediately using "/usr/bin/udisksctl mount" after this script.
sudo partprobe $SDCARD
sync
sleep 1
while [ ! -b ${SDCARD}2 ]; do
  sleep 1
done
