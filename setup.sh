#!/bin/sh
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Install/update everything needed to build and flash this project on a Ubuntu
# host.
#
# This script doesn't need root.
#
# This is based on instructions are:
# - https://github.com/pfalcon/esp-open-sdk
# - https://github.com/SmingHub/Sming/wiki/Linux-Quickstart


set -eu


echo "Please make sure prerequisites are installed:"
echo "sudo apt-get install \\"
echo "    autoconf automake bash bc bison flex g++ gawk gcc git gperf \\"
echo "    help2man libexpat1-dev libncurses5-dev libtool make python \\"
echo "    python-serial sed srecord texinfo unrar unzip"
echo ""


echo "- Processing esp-open-sdk. This is by far the slowest step"
# Commit: 90eb4a8d833e7595282178e832121351ab6f3b90
if [ -d "$ESP_HOME" ]; then
  echo "  Pulling"
  cd "$ESP_HOME"
  git pull
else
  echo "  Checking out"
  git clone --recursive https://github.com/pfalcon/esp-open-sdk "$ESP_HOME"
  cd "$ESP_HOME"
fi
echo "- Building"
make
echo ""


echo "- Installing esptool"
echo ""
# Remove --user to install system wide.
pip install --user --upgrade esptool
echo ""


echo "- Checking out and build esptool2"
# Commit: ec0e2c72952f4fa8242eedd307c58a479d845abe
if [ -d "$ESP_HOME" ]; then
  echo "  Pulling"
  cd "$ESP_HOME/../esptool2"
  git pull
else
  echo "  Checking out"
  git clone https://github.com/raburton/esptool2 "$ESP_HOME/../esptool2"
  cd "$ESP_HOME/../esptool2"
fi
make
echo ""


echo "- Checking out and building Sming"
echo ""
if [ -d "$SMING_HOME/.." ]; then
  echo "  Pulling"
  cd "$SMING_HOME/.."
  git pull
else
  echo "  Checking out"
  # Append for the stable branch: -b master
  git clone https://github.com/SmingHub/Sming "$SMING_HOME/.."
fi
cd "$SMING_HOME"
make
make spiffy
