#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Change locale to Canadian.

set -eu

# Use "timedatectl list-timezones" to list the values.
timedatectl set-timezone America/Toronto

if [ "$(grep 'ID=' /etc/os-release)" == "ID=raspbian" ]; then
  # Use the us keyboard layout.
  sed -i 's/XKBLAYOUT="gb"/XKBLAYOUT="us"/' /etc/default/keyboard
  # Fix Wifi country settings for Canada.
  raspi-config nonint do_wifi_country CA

  # Switch to en_US.
  sed -i 's/en_GB/en_US/' /etc/locale.gen
  dpkg-reconfigure --frontend=noninteractive locales
  update-locale LANG=en_US.UTF-8

  #sed -i s/en_GB/en_US/ /etc/default/locale
  #sed -i -e "s/# $LANG.*/$LANG.UTF-8 UTF-8/" /etc/locale.gen
  #locale-gen --purge en_US.UTF-8
  #echo -e 'LANG="en_US.UTF-8"\nLANGUAGE="en_US:en"\n' > /etc/default/locale
fi
