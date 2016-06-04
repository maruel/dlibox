#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# This script is useful to run the tool continuously without installing it as a
# systemd service.

# Use iptables to redirect incoming TCP on port 80 to port 8080:
#
#   sudo iptables -A PREROUTING -t nat -p tcp --dport 80 -j REDIRECT --to-port 8080
#   sudo iptables -t nat -A OUTPUT -o lo -p tcp --dport 80 -j REDIRECT --to-port 8080

set -u

# trap ctrl-c and call ctrl_c()
trap ctrl_c INT

function ctrl_c() {
  echo "** Trapped CTRL-C"
  exit 1
}

while true; do
  echo "$(date --rfc-3339=seconds): $GOPATH/bin/dlibox -verbose -port 8080"
  $GOPATH/bin/dlibox -verbose -port 8080
done
