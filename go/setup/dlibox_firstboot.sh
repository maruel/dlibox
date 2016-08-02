#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

if [ "$USER" != "root" ]; then
  echo "raspbian.sh can only be run as root"
  exit 1
fi

LOG_FILE=/var/log/dlibox_firstboot.log
if [ -f $LOG_FILE ]; then
  exit 1
fi

# Close stdout and stderr
exec 1<&-
exec 2<&-

# Open stdout as $LOG_FILE file for read and write.
exec 1<>$LOG_FILE
# Redirect stderr to stdout.
exec 2>&1

echo "- Configuring SSH, SPI, I2C"
sed -i 's/PasswordAuthentication yes/#PasswordAuthentication no/' /etc/ssh/sshd_config
# https://github.com/RPi-Distro/raspi-config/blob/master/raspi-config
raspi-config nonint do_ssh 0
raspi-config nonint do_spi 0
raspi-config nonint do_i2c 0

echo "- Updating OS"
apt-get update
apt-get upgrade -y
apt-get install -y git ifstat ntpdate sysstat tmux vim

echo "- Installing as user"
sudo -u pi -- <<EOF
mkdir /home/pi/bin
cd /home/pi; git clone --recurse https://github.com/maruel/bin_pub /home/pi/bin/bin_pub
/home/pi/bin/bin_pub/setup/update_config.py
/home/pi/bin/bin_pub/setup/install_golang.py
go get github.com/maruel/dlibox/go/cmd/dlibox
EOF

# TODO(maruel): Change locale & keyboard layout (?)
# dpkg-reconfigure locales
# TODO(maruel): Change timezone.
# dpkg-reconfigure tzdata
# raspi-config nonint do_wifi_country CA
# TODO(maruel): Change hostname to unique name.
# raspi-config nonint do_hostname dlibox-<serial>
# TODO(maruel): Reduce GPU memory.
# TODO(maruel): avahi browse to detect if MQTT is installed, install otherwise.
# TODO(maruel): Change 'pi' password?

echo "- Done, rebooting"
reboot
#echo "- Extending partition"
#raspi-config --expand-rootfs
