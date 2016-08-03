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

if [ -d ${SDCARD}1 ]; then
  echo "- Unmounting"
  for i in ${SDCARD}?; do
    echo "  $i"
    umount $i || true
  done
fi

if [ ! -f $IMGNAME.img ]; then
  echo "- Fetching Raspbian Jessie Lite latest"
  curl -LO https://downloads.raspberrypi.org/raspbian_lite_latest
  unzip $IMGNAME.zip
  rm $IMGNAME.zip
fi

echo "- Flashing (takes 2 minutes)"
sudo time dd bs=4M if=$IMGNAME.img of=$SDCARD


mkdir mounted_img


# Skip this if you don't use a small display.
if [ true ]; then
  # Mount the SDCard do to small modifications. Do not modify the original
  # image.
  echo "- Mounting 'boot' partition on the SDCard for modification"
  sudo mount ${SDCARD}1 mounted_img

  # Strictly speaking, you won't need a monitor at all since ssh will be up and
  # running and the device will connect to the SSID provided.
  # Search for [5 Inch 800x480], found one at 23$USD with fre shipping on
  # aliexpress.
  echo "- Enabling 5\" display support (optional)"
  sudo tee --append mounted_img/config.txt > /dev/null <<EOF

# Enable support for 800x480 display:
hdmi_group=2
hdmi_mode=87
hdmi_cvt 800 480 60 6 0 0 0

# Enable touchscreen:
# Not necessary on Jessie Lite since it boots in console mode. :)
# Some displays use 22, others 25.
#dtoverlay=ads7846,penirq=22,penirq_pull=2,speed=10000,xohms=150

EOF

  echo "- Unmounting boot"
  sudo umount mounted_img
fi


# Setup SSH keys, wifi and automatic setup process.
if [ true ]; then
  # TODO(maruel): Formatting to F2FS would be nice but this requires one boot on
  # the rPi to be able to run "apt-get install f2fs-tools" first. I don't know
  # how to do it otherwise.
  # http://whitehorseplanet.org/gate/topics/documentation/public/howto_ext4_to_f2fs_root_partition_raspi.html

  echo "- Mounting 'root' partition on the SDCard for modification"
  sudo mount ${SDCARD}2 mounted_img

  echo "- Copying dlibox_firstboot.sh"
  sudo cp dlibox_firstboot.sh mounted_img/etc/init.d
  echo "- Copying ~/.ssh/authorized_keys"
  sudo mkdir mounted_img/home/pi/.ssh
  # This assumes you have properly set your own ssh keys and plan to use them.
  sudo cp $HOME/.ssh/authorized_keys mounted_img/home/pi/.ssh/authorized_keys
  # pi(1000).
  sudo chown -R 1000:1000 mounted_img/home/pi/.ssh
  # Otherwise network cable works.
  echo "- Setting up wifi"
  sudo tee --append mounted_img/etc/wpa_supplicant/wpa_supplicant.conf > /dev/null <<EOF

network={
  ssid="$SSID"
  psk="$WIFI_PWD"
}
EOF

  echo "- Unmounting root"
  sudo umount mounted_img
fi


rm -r mounted_img
sync
