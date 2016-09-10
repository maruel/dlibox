#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# TODO(someone): Make this script OSX compatible. For now it was only tested on
# Ubuntu.

set -eu

if [ "$#" -ne 1 ]; then
  echo "Unmount all partitions on a SD card"
  echo ""
  echo "usage: ./umount.sh /dev/<sdcard_path>"
  exit 1
fi


for i in $1*; do
  if [ "$i" != "$1" ]; then
    /usr/bin/udisksctl unmount -f -b $i || true
  fi
done
