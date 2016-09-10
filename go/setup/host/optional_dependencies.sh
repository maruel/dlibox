#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Install *optional* dlibox Armbian and Raspbian dependencies.

set -eu
cd "$(dirname $0)"


## MQTT server

apt-get -y install mosquitto
cp mosquitto.conf /etc/mosquitto/conf.d/
systemctl restart mosquitto


### OpenHAB

# OpenHAB integrates ntp and mdns so maybe disable avahi and remove ntpdate?

# You need:
# https://github.com/eclipse/smarthome/blob/master/docs/documentation/community/downloads.md#designer-builds
# or
# http://www.openhab.org/getting-started/downloads.html

# http://docs.openhab.org/installation/linux.html#package-repository-installation
# http://brabraen.blogspot.ca/2016/04/pine64-install-oracle-jdk-8.html
# http://www.oracle.com/technetwork/java/javase/downloads/java-archive-javase8-2177648.html
# OpenHAB2 requires Java 8. Raspbian includes it but not Armbian.
#echo 'deb https://httpredir.debian.org/debian/ jessie main contrib' > /etc/apt/sources.list.d/java.list

# Install the OpenHab repo so it will be kept up to date. This includes systemd
# files.
wget -qO - 'https://bintray.com/user/downloadSubjectPublicKey?username=openhab' | apt-key add -
echo 'deb http://dl.bintray.com/openhab/apt-repo2 testing main' > /etc/apt/sources.list.d/openhab2.list
apt-get update
apt-get -y install oracle-java8-jdk openhab2-online

systemctl daemon-reload
systemctl enable openhab2.service
systemctl start openhab2.service

# HTTP Binding
# MQTT Binding
# Network Binding

# Configuration files are at:
# /etc/openhab2/*
# /var/lib/openhab2/*
# /var/log/openhab2/*


# Samba

apt-get install -y samba samba-common-bin
# echo 'wins support = yes' >> /etc/samba/smb.conf
echo >> /etc/samba/smb.conf  <<'EOF'
[openHAB-sys]
  comment=openHAB2 application
  path=/usr/share/openhab2
  browseable=Yes
  writeable=Yes
  only guest=no
  public=no
  create mask=0777
  directory mask=0777

[openHAB-conf]
  comment=openHAB2 site configuration
  path=/etc/openhab2
  browseable=Yes
  writeable=Yes
  only guest=no
  public=no
  create mask=0777
  directory mask=0777
EOF

# smbpasswd -a openhab
systemctl restart smbd.service
