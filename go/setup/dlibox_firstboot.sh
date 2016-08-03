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


echo "- Configuring SSH, SPI, I2C, GPU memory"
# Force key based authentication since the password is known.
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
# https://github.com/RPi-Distro/raspi-config/blob/master/raspi-config
# 0 means enabled.
raspi-config nonint do_ssh 0
raspi-config nonint do_spi 0
raspi-config nonint do_i2c 0
# Lowers GPU memory from 64Mb to 16Mb. Doing so means goodbye to startx.
raspi-config nonint do_memory_split 16


echo "- Configuring locale as French Canadian"
# Use the best keyboard layout. Change to "us" if needed. :)
sed -i s/XKBLAYOUT="gb"/XKBLAYOUT="ca"/ /etc/default/keyboard
# Use "timedatectl list-timezones" to list the values.
timedatectl set-timezone America/Toronto
# Switch to en_US.
locale-gen --purge en_US.UTF-8
sed -i s/en_GB/en_US/ /etc/default/locale
# Fix Wifi country settings for Canda.
raspi-config nonint do_wifi_country CA


echo "- Updating OS"
apt-get update
apt-get upgrade -y
apt-get install -y git ifstat ntpdate sysstat tmux vim


# TODO(maruel): Do not forget to update the Go version as needed.
echo "- Installing as user"
sudo -u pi -- <<'EOF'
cd
mkdir src
git clone --recurse https://github.com/maruel/bin_pub bin/bin_pub
bin/bin_pub/setup_scripts/update_config.py
export GOROOT=/home/pi/src/golang
curl -S https://storage.googleapis.com/golang/go1.6.3.linux-armv6l.tar.gz | tar xz
mv go src/golang
#bin/bin_pub/setup/install_golang.py
PATH="$PATH:$GOROOT/bin"
export GOPATH="$HOME/src/gopath"
go get github.com/maruel/dlibox/go/cmd/dlibox
EOF

# TODO(maruel): Setup dlibox.service


# TODO(maruel): Change hostname to unique name. The problem is that it becomes
# harder to find the device on the network.
#SERIAL=$(cat /proc/cpuinfo | grep Serial | cut -d ":" -f 2 | sed 's/^ 0\+//')
# raspi-config nonint do_hostname dlibox-$SERIAL

# TODO(maruel): avahi browse to detect if MQTT is installed, install otherwise.
# TODO(maruel): kbdrate -d 200 -r 60

echo "- Done, rebooting"
reboot
#echo "- Extending partition"
#raspi-config --expand-rootfs
