#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Change hostname to unique name. The problem is that it becomes harder to find
# the device on the network but this is necessary when configuring multiple
# devices. To find it, use ../find.sh.
#
# Since the hostname is based on the serial number of the CPU with leading zeros
# trimmed off, it is a constant yet unique value.

set -eu

SERIAL=$(cat /proc/cpuinfo | grep Serial | cut -d ":" -f 2 | sed 's/^ 0\+//')
HOST=dlibox-$SERIAL
echo "New hostname is: $HOST"
raspi-config nonint do_hostname $HOST
