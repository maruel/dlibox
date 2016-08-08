#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Change locale to Canadian.

set -eu

# Use the us keyboard layout.
sed -i 's/XKBLAYOUT="gb"/XKBLAYOUT="us"/' /etc/default/keyboard
# Use "timedatectl list-timezones" to list the values.
timedatectl set-timezone America/Toronto
# Switch to en_US.
locale-gen --purge en_US.UTF-8
sed -i s/en_GB/en_US/ /etc/default/locale
# Fix Wifi country settings for Canada.
raspi-config nonint do_wifi_country CA
