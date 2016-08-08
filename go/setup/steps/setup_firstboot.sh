#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.


# Modifies the root partition of Raspbian:
# - setup wifi connectivity.
# - add ssh key for passwordless login.
# - setup dlibox_firstboot.sh to start automatically, to complete the setup.
#
# Note:
# dlibox_firstboot.sh will eventually disable password based ssh login, since
# the password is a constant. But there is a time window of a few minutes where
# it's possible to login with pi/raspberry credentials. If the network you are
# using is potentially unsafe, it could be worth changing /etc/ssh/sshd_config
# here.

set -eu
cd "$(dirname $0)"


if [ "$#" -ne 2 ]; then
  echo "usage: ./setup_firstboot.sh <root partition> <ssid>"
  exit 1
fi


cd "$(dirname $0)"


ROOT_PATH=$1
if [ ! -d $ROOT_PATH/home/pi ]; then
  echo "If not running on a Raspberry Pi, provide the path to the mounted partition"
  exit 1
fi


# TODO(maruel): Formatting to F2FS would be nice but this requires one boot on
# the rPi to be able to run "apt-get install f2fs-tools" first. I don't know
# how to do it otherwise.
# http://whitehorseplanet.org/gate/topics/documentation/public/howto_ext4_to_f2fs_root_partition_raspi.html


echo "- First boot setup script"
sudo cp ../support/dlibox_firstboot.sh $ROOT_PATH/root
sudo chmod +x $ROOT_PATH/root/dlibox_firstboot.sh
# Skip this step to debug dlibox_firstboot.sh. Then login at the console and run
# the script manually.
sudo mv $ROOT_PATH/etc/rc.local $ROOT_PATH/etc/rc.local.old
sudo tee $ROOT_PATH/etc/rc.local > /dev/null <<'EOF'
#!/bin/sh -e
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# As part of https://github.com/maruel/dlibox

LOG_FILE=/var/log/dlibox_firstboot.log
if [ ! -f $LOG_FILE ]; then
  /root/dlibox_firstboot.sh 2>&1 | tee $LOG_FILE
fi
exit 0
EOF
sudo chmod +x $ROOT_PATH/etc/rc.local


echo "- SSH keys"
# This assumes you have properly set your own ssh keys and plan to use them.
sudo mkdir $ROOT_PATH/home/pi/.ssh
sudo cp $HOME/.ssh/authorized_keys $ROOT_PATH/home/pi/.ssh/authorized_keys
# pi(1000).
sudo chown -R 1000:1000 $ROOT_PATH/home/pi/.ssh
# Force key based authentication since the password is known.
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' $ROOT_PATH/etc/ssh/sshd_config


echo "- Wifi"
# TODO(maruel): Get the data from /etc/NetworkManager/system-connections/*
SSID="$2"
# TODO(maruel): When not found, ask the user for the password. It's annoying to
# test since the file is only readable by root.
# TODO(maruel): Ensure it works with SSID with whitespace/emoji in their name.
WIFI_PWD="$(sudo grep -oP '(?<=psk=).+' /etc/NetworkManager/system-connections/$SSID)"
sudo tee --append $ROOT_PATH/etc/wpa_supplicant/wpa_supplicant.conf > /dev/null <<EOF

network={
  ssid="$SSID"
  psk="$WIFI_PWD"
}
EOF
