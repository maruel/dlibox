#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Make apt-get update & upgrade run on a scheduled basis.

set -eu
cd "$(dirname $0)"

HOST="$(./gen_hostname.sh)"

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
