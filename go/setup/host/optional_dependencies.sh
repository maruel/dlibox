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
