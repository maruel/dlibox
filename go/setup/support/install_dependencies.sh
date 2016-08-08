#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Install dlibox Raspbian dependencies.
# Change hardware settings to enable I²C, SPI and lower GPU memory usage.

set -eu
cd "$(dirname $0)"

# https://github.com/RPi-Distro/raspi-config/blob/master/raspi-config
# 0 means enabled.
raspi-config nonint do_spi 0
raspi-config nonint do_i2c 0
# Lowers GPU memory from 64Mb to 16Mb. Doing so means goodbye to startx.
raspi-config nonint do_memory_split 16

# lirc

apt-get -y install lirc

# Configure lirc to use the lir-rpi driver including in /boot/overlays.
sed -i s'/DRIVER="UNCONFIGURED"/DRIVER="default"/' /etc/lirc/hardware.conf
sed -i s'/DEVICE=""/DEVICE="\/dev\/lirc0"/' /etc/lirc/hardware.conf
sed -i s'/MODULES=""/MODULES="lirc_rpi"/' /etc/lirc/hardware.conf

mv /etc/lirc/lircd.conf /etc/lirc/lircd.conf.org
cp lircd.conf /etc/lirc/

# TODO(maruel): Do not add twice.
echo "" >> /boot/config.txt
echo "# dlibox" >> /boot/config.txt
echo "# https://github.com/raspberrypi/firmware/blob/master/boot/overlays/README" >> /boot/config.txt
echo "dtoverlay=lirc-rpi,gpio_in_pull=high" >> /boot/config.txt
