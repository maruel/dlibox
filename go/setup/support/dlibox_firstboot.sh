#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# As part of https://github.com/maruel/dlibox

set -eu


echo "- Installing git"
# apt-get update must be done right away, since the old packages are likely not
# on the mirrors anymore.
# Defer apt-get upgrade for later. It'll be done via a nightly cron job. Doing
# it now takes several minutes and the user is eagerly waiting for the lights to
# come up!
apt-get update
apt-get install -y git


echo "- Installing as user - setting up GOPATH, GOROOT and PATH"
cat >> /home/pi/.profile <<'EOF'

export GOROOT=$HOME/go
export GOPATH=$HOME
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

EOF

# TODO(maruel): Do not forget to update the Go version as needed.
echo "- Installing as user - Go and dlibox"
sudo -i -u pi /bin/sh <<'EOF'
cd
curl -S https://storage.googleapis.com/golang/go1.6.3.linux-armv6l.tar.gz | tar xz
go get -v github.com/maruel/dlibox/go/cmd/dlibox
EOF

# At this point, defer control to the script in the repository.
/home/pi/src/github.com/maruel/dlibox/go/setup/support/finish_install.sh
