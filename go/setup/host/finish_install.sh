#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Follow up from dlibox_firstboot.sh
# You can run it manually if you did not want to use the scripts to generate the
# modified Raspbian image.
# This script assumes Go was already installed and that dlibox was installed
# via:
#  go get github.com/maruel/dlibox/go/cmd/dlibox

set -eu
cd "$(dirname $0)"

if [ "${USER:=root}" != "root" ]; then
  echo "This script must be run as root."
  exit 1
fi

echo "- Changing hostname"
./set_hostname.sh

echo "- Configuring locale as Canadian"
./set_locale.sh

echo "- Installing dependencies and Configuring SPI, I2C, GPU memory"
./install_dependencies.sh

echo "- Setting up dlibox as a service and auto-update cron job"
./install_systemd.sh

echo "- Setting up automated apt cron job"
./schedule_apt.sh

echo "- Installing ancillary utilities"
sudo -i -u pi /bin/sh <<'EOF'
go get -v github.com/maruel/dlibox/go/cmd/... github.com/maruel/dlibox/go/pio/cmd/...
EOF

echo "- User specific installation"
./user_config.sh

# TODO(maruel): Configure outgoing email for notifications. Probably in
# user_config.sh?

echo "- Done, rebooting"
reboot
