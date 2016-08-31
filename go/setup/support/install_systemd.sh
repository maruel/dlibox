#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Install dlibox as a service.
# Do not start it.

set -eu
cd "$(dirname $0)"

cp ../systemd/* /etc/systemd/system
systemctl daemon-reload
systemctl enable dlibox.service
#systemctl enable dlibox_update.service
systemctl enable dlibox_update.timer
