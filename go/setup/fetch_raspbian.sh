#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

IMGNAME=2016-05-27-raspbian-jessie-lite

if [ "$#" -ne 3 ]; then
  echo "usage: ./fetch_raspbian.sh /dev/<sdcard_path> <ssid> <wifi_pwd>"
  exit 1
fi

SDCARD=$1
SSID="$2"
WIFI_PWD="$3"

echo "- Unmounting"
for i in ${SDCARD}?; do
  echo "  $i"
  umount $i || true
done

if [ ! -f $IMGNAME.img ]; then
  echo "- Fetching Raspbian Jessie Lite latest"
  curl -LO https://downloads.raspberrypi.org/raspbian_lite_latest
  unzip $IMGNAME.zip
  rm $IMGNAME.zip
fi

echo "- Flashing"
sudo dd bs=4M if=$IMGNAME.img of=$SDCARD

# TODO(maruel): Formatting to F2FS would be nice but this requires one boot on
# the rPi to be able to run "apt-get install f2fs-tools" first. I don't know how
# to do it otherwise.
# http://whitehorseplanet.org/gate/topics/documentation/public/howto_ext4_to_f2fs_root_partition_raspi.html

echo "- Mounting boot"
mkdir foo
sudo mount ${SDCARD}1 foo

echo "- Enabling 5\" display"
sudo tee --append foo/config.txt > /dev/null <<EOF

# Enable support for 800x480 display:
hdmi_group=2
hdmi_mode=87
hdmi_cvt 800 480 60 6 0 0 0

# Enable touchscreen, not necessary on Jessie Lite since it boots in console
# mode.
# Some displays use 22, others 25.
dtoverlay=ads7846,penirq=22,penirq_pull=2,speed=10000,xohms=150
EOF

echo "- Unmounting boot"
sudo umount foo

echo "- Mounting root"
sudo mount ${SDCARD}2 foo

echo "- Copying dlibox_firstboot.sh"
sudo cp dlibox_firstboot.sh foo/etc/init.d
echo "- Copying ~/.ssh/authorized_keys"
sudo mkdir foo/home/pi/.ssh
sudo cp $HOME/.ssh/authorized_keys foo/home/pi/.ssh/authorized_keys
# TODO(maruel): chown with the right user id.
echo "- Setting up wifi"
sudo tee --append foo/etc/wpa_supplicant/wpa_supplicant.conf > /dev/null <<EOF

network={
  ssid="$SSID"
  psk="$WIFI_PWD"
}
EOF

echo "- Unmounting root"
sudo umount foo
rm -r foo

sync
