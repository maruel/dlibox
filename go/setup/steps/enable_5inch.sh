#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Search for [5 Inch 800x480], found one at 23$USD with free shipping on
# aliexpress.

set -eu
cd "$(dirname $0)"


if [ "$#" -ne 1 ]; then
  echo "usage: ./enable_5inch.sh <path to /boot>"
  exit 1
fi


BOOT_PATH=$1


if [ ! -f $BOOT_PATH/config.txt ]; then
  echo "usage: ./enable_5inch.sh <path to /boot>"
  exit 1
fi


sudo tee --append $BOOT_PATH/config.txt > /dev/null <<EOF

# Enable support for 800x480 display:
hdmi_group=2
hdmi_mode=87
hdmi_cvt 800 480 60 6 0 0 0

# Enable touchscreen:
# Not necessary on Jessie Lite since it boots in console mode. :)
# Some displays use 22, others 25.
# Enabling this means the SPI bus cannot be used anymore.
#dtoverlay=ads7846,penirq=22,penirq_pull=2,speed=10000,xohms=150

EOF
