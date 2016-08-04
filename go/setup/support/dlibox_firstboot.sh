#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# As part of https://github.com/maruel/dlibox

set -eu


# When run via /etc/rc.local, USER is not defined.
if [ "${USER:=root}" != "root" ]; then
  echo "OMG don't run this locally!"
  exit 1
fi


# The idea is that this command will fail if not running on Raspbian, as a
# safety measure.
echo "- Testing if running on Raspbian"
grep raspbian /etc/os-release > /dev/null


echo "- Changing hostname"
# Change hostname to unique name. The problem is that it becomes harder to find
# the device on the network but this is necessary when configuring multiple
# devices. Hint: find it with:
#   avahi-browse -t _workstation._tcp -l -k | grep IPv4
#
# Since the hostname is based on the serial number of the CPU with leading zeros
# trimmed off, it is a constant yet unique value.
SERIAL=$(cat /proc/cpuinfo | grep Serial | cut -d ":" -f 2 | sed 's/^ 0\+//')
HOST=dlibox-$SERIAL
echo "  New hostname is: $HOST"
raspi-config nonint do_hostname $HOST


echo "- Configuring SSH, SPI, I2C, GPU memory"
# Force key based authentication since the password is known.
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
# https://github.com/RPi-Distro/raspi-config/blob/master/raspi-config
# 0 means enabled.
raspi-config nonint do_spi 0
raspi-config nonint do_i2c 0
# Lowers GPU memory from 64Mb to 16Mb. Doing so means goodbye to startx.
raspi-config nonint do_memory_split 16


echo "- Configuring locale as Canadian"
# Use the us keyboard layout.
sed -i 's/XKBLAYOUT="gb"/XKBLAYOUT="us"/' /etc/default/keyboard
# Use "timedatectl list-timezones" to list the values.
timedatectl set-timezone America/Toronto
# Switch to en_US.
locale-gen --purge en_US.UTF-8
sed -i s/en_GB/en_US/ /etc/default/locale
# Fix Wifi country settings for Canada.
raspi-config nonint do_wifi_country CA


echo "- Updating OS"
apt-get update
apt-get upgrade -y
apt-get install -y git ifstat ntpdate sysstat tmux vim
apt-get autoclean


echo "- Installing as user"
cat >> /home/pi/.profile <<'EOF'

export GOROOT=$HOME/go
export GOPATH=$HOME
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

EOF

# TODO(maruel): Do not forget to update the Go version as needed.
# Running bin/bin_pub/setup/install_golang.py would unconditionally install the
# latest version but it is slower to run (several minutes) than just fetching a
# known good version.
sudo -i -u pi /bin/sh <<'EOF'
cd
git clone --recurse https://github.com/maruel/bin_pub bin/bin_pub
#bin/bin_pub/setup_scripts/update_config.py
curl -S https://storage.googleapis.com/golang/go1.6.3.linux-armv6l.tar.gz | tar xz
go get -v github.com/maruel/dlibox/go/cmd/dlibox
EOF


echo "- Setting up dlibox as a service and auto-update timer"
# Copy and enable the 2 services but do not start them, the host will soon
# reboot.
cp /home/pi/src/github.com/maruel/dlibox/go/setup/support/dlibox.service /etc/systemd/system
cp /home/pi/src/github.com/maruel/dlibox/go/setup/support/dlibox_update.service /etc/systemd/system
cp /home/pi/src/github.com/maruel/dlibox/go/setup/support/dlibox_update.timer /etc/systemd/system
systemctl daemon-reload
systemctl enable dlibox.service
#systemctl enable dlibox_update.service
systemctl enable dlibox_update.timer


echo "- Setting up automated upgrade"
# This runs through /etc/cron.daily/apt. More details at
# https://wiki.debian.org/UnattendedUpgrades
# TODO(maruel): Confirm the following to work.
# TODO(maruel): Does this need apt-get install unattended-upgrades
# update-notifier-common ?
# TODO(maruel): Configure /etc/apt/listchanges.conf
echo > /etc/apt/apt.conf.d/30dlibox <<EOF
# Enable /etc/cron.daily/apt.
APT::Periodic::Enable "1";
# apt-get update every 7 days.
APT::Periodic::Update-Package-Lists "7";
# apt-get upgrade every 7 days.
APT::Periodic::Unattended-Upgrade "7";
# apt-get autoclean evey 21 days.
APT::Periodic::AutocleanInterval "21";
# Log a bit more.
APT::Periodic::Verbose "1";
# automatically reboot after updates.
Unattended-Upgrade::Automatic-Reboot "true";
# apt-get autoremove.
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Mail "root@$HOST";
EOF


# TODO(maruel): Configure outgoing email for notifications.


echo "- Done, rebooting"
reboot
