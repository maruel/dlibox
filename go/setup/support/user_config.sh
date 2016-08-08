#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Changes that are specific to me (Marc-Antoine Ruel), that I expect you to want
# to change or comment out.

set -eu

apt-get install -y ifstat ntpdate sysstat tmux vim

sudo -i -u pi /bin/sh <<'EOF'
cd
git clone --recurse https://github.com/maruel/bin_pub bin/bin_pub
bin/bin_pub/setup_scripts/update_config.py
EOF
