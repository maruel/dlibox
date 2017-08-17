#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# As part of https://github.com/maruel/dlibox

set -eu


function install_lirc() {
  echo "- Setting up lirc"
  sudo apt-get -y install lirc

  # Configure lirc to use the lir-rpi driver including in /boot/overlays.
  sudo sed -i s'/DRIVER="UNCONFIGURED"/DRIVER="default"/' /etc/lirc/hardware.conf
  sudo sed -i s'/DEVICE=""/DEVICE="\/dev\/lirc0"/' /etc/lirc/hardware.conf
  sudo sed -i s'/MODULES=""/MODULES="lirc_rpi"/' /etc/lirc/hardware.conf

  sudo mv /etc/lirc/lircd.conf /etc/lirc/lircd.conf.org
  sudo tee etc/lirc/lircd.conf > /dev/null <<EOF
# https://github.com/maruel/dlibox
begin remote
  name            ElectroDragon
  bits            16
  flags           SPACE_ENC|CONST_LENGTH
  eps             30
  aeps            100
  header          9000  4500
  one             563  1687
  zero            563  562
  ptrail          563
  repeat          9000  2250
  pre_data_bits   16
  pre_data        0xFF
  gap             108000
  toggle_bit_mask 0x0
  frequency       38000
  duty_cycle      33
  begin codes
      KEY_CHANNELDOWN          0xA25D
      KEY_CHANNEL              0x629D
      KEY_CHANNELUP            0xE21D
      KEY_PREVIOUS             0x22DD
      KEY_NEXT                 0x02FD
      KEY_PLAYPAUSE            0xC23D
      KEY_VOLUMEUP             0xA857
      KEY_VOLUMEDOWN           0xE01F
      KEY_EQ                   0x00FF906F
      KEY_NUMERIC_0            0x6897
      KEY_100PLUS              0x9867
      KEY_200PLUS              0xB04F
      KEY_NUMERIC_1            0x30CF
      KEY_NUMERIC_2            0x18E7
      KEY_NUMERIC_3            0x7A85
      KEY_NUMERIC_4            0x10EF
      KEY_NUMERIC_5            0x38C7
      KEY_NUMERIC_6            0x5AA5
      KEY_NUMERIC_7            0x42BD
      KEY_NUMERIC_8            0x4AB5
      KEY_NUMERIC_9            0x52AD
  end codes
end remote
EOF

  # TODO(maruel): Do not add twice.
  # TODO(maruel): Only on Raspbian.
  sudo tee --append /boot/config.txt > /dev/null <<EOF

# dlibox
# https://github.com/raspberrypi/firmware/blob/master/boot/overlays/README
dtoverlay=lirc-rpi,gpio_out_pin=5,gpio_in_pin=13,gpio_in_pull=high

# dtoverlay=spi1-1cs
dtoverlay=pi3-disable-bt
EOF
}


function install_mqtt() {
  echo "- Installing MQTT server"
  sudo apt-get -y install mosquitto
  sudo tee /etc/mosquitto/conf.d/mosquitto.conf > /dev/null <<EOF
# https://github.com/maruel/dlibox
# TODO(maruel): Change!
allow_anonymous true

listener 1883
listener 1884
protocol websockets

# General settings. Assumes Debian default mosquitto.conf was used.
allow_duplicate_messages false
message_size_limit 65536
retry_interval 5

# Options to use potentially.
# TODO(maruel): Create self-signed certificate and distribute to the devices.
# For ESP8266, use library that supports certificates, e.g.
# https://github.com/tuanpmt/esp_mqtt
# certfile file path
# tls_version tlsv1.2
# acl_file /etc/mosquitto/acl.txt
# password_file /etc/mosquitto/pwd.txt
# psk_file /etc/mosquitto/psk.txt
EOF
  sudo systemctl restart mosquitto
}


function install_dlibox() {
  AS_USER=$1

  echo "- Injecting history in .bash_history"
  cat >> /home/${AS_USER}/.bash_history <<'EOF'
sudo systemctl stop dlibox
sudo journalctl -f -u dlibox
EOF

  echo "- Installing dlibox as user {$AS_USER}"
  go get -u -v github.com/maruel/dlibox/cmd/dlibox

  echo "- Setting up dlibox as system service"
  sudo tee /etc/systemd/system/dlibox.service > /dev/null <<EOF
# https://github.com/maruel/dlibox
[Unit]
Description=Runs dlibox automatically upon boot
Wants=network-online.target
After=network-online.target

[Service]
User=${AS_USER}
Group=${AS_USER}
KillMode=mixed
Restart=always
TimeoutStopSec=20s
ExecStart=/home/${AS_USER}/go/bin/dlibox
# Systemd 229:
#AmbientCapabilities=CAP_NET_BIND_SERVICE
# Systemd 228 and below:
SecureBits=keep-caps
Capabilities=cap_net_bind_service+pie
# Older systemd:
PermissionsStartOnly=true
ExecStartPre=/sbin/setcap 'cap_net_bind_service=+ep' /home/${AS_USER}/go/bin/dlibox
# High priority stuff:
# Nice=-20
# IOSchedulingClass=realtime
# IOSchedulingPriority=0
# CPUSchedulingPolicy=rr
# CPUSchedulingPriority=99
# CPUSchedulingResetOnFork=true

[Install]
WantedBy=default.target
EOF

  sudo tee /etc/systemd/system/dlibox_update.service > /dev/null <<EOF
# https://github.com/maruel/dlibox
[Unit]
Description=Updates dlibox, as triggered by dlibox_update.timer
After=network-online.target
[Service]
Type=oneshot
User=${AS_USER}
Group=${AS_USER}
# /bin/sh is necessary to load .profile to set $GOPATH:
ExecStart=/bin/sh -l -c "go get -v -u github.com/maruel/dlibox/cmd/dlibox"
WorkingDirectory=/home/${AS_USER}
EOF

  # TODO(maruel): Nightly cron job at 4:18.
  sudo tee /etc/systemd/system/dlibox_update.timer > /dev/null <<EOF
[Unit]
Description=go get -u dlibox as a cron job
[Timer]
OnBootSec=1min
OnUnitActiveSec=10min
RandomizedDelaySec=5
[Install]
WantedBy=timers.target
EOF

  sudo systemctl daemon-reload
  sudo systemctl enable dlibox.service
  sudo systemctl enable dlibox_update.timer
}


function install_optional() {
  echo "- Installing optional tools"
  git clone --recurse https://github.com/maruel/bin_pub bin/bin_pub
  bin/bin_pub/setup_scripts/update_config.py
  go get -v \
    github.com/FiloSottile/gorebuild \
    github.com/maruel/panicparse/cmd/pp \
    golang.org/x/tools/cmd/goimports
}


function do_controller() {
  install_mqtt
  install_dlibox $USER
  #install_optional
}


function do_device() {
  install_lirc
  install_dlibox $USER
  #install_optional
}


#do_controller
#do_device
