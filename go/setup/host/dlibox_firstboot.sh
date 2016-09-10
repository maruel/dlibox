#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# As part of https://github.com/maruel/dlibox

set -eu


echo "- Injecting history in .bash_history"
cat >> /home/pi/.bash_history <<'EOF'
sudo systemctl stop dlibox
sudo journalctl -f -u dlibox
tail -f /var/log/dlibox_firstboot.log
EOF


echo "- Installing git"
# apt-get update must be done right away, since the old packages are likely not
# on the mirrors anymore.
# Defer apt-get upgrade for later. It'll be done via a nightly cron job. Doing
# it now takes several minutes and the user is eagerly waiting for the lights to
# come up!
apt-get update
apt-get install -y git


if [ "$(getconf LONG_BIT)" == "64" ]; then
  # This is more complex as Armbian on A64 doesn't have 32 bit userland support
  # installed by default, and go1.4.3 cannot compile natively on aarch64
  # (arm64). So install golang 1.6 via apt-get then use it to bootstrap a more
  # recent (1.7 at the time of writting) version.
  # That doesn't work because by default armbian doesn't have a swap file,
  # causing https://github.com/golang/go/issues/16082
  echo "- Installing as user - Go"
  apt-get install -y golang
#  sudo -i -u pi /bin/sh <<'EOF'
#git clone https://go.googlesource.com/go ~/go
#cd ~/go/src
#git checkout "$(git tag | grep "^go" | egrep -v "beta|rc" | tail -n 1)"
#GOROOT_BOOTSTRAP=/usr/lib/go ./make.bash
#EOF
#  # Then remove the outdated system version.
#  apt-get remove -y golang
#  # Then copy itself to a backup so the next upgrade is seamless.
#  cp -a /home/pi/go /home/pi/go1.4
fi
  # Fast path to skip compilation on 32 bits hosts. This uses the slower armv6l
  # binaries.
  # TODO(maruel): Do not forget to update the Go version as needed.
  echo "- Installing as user - setting up GOROOT"
cat >> /home/pi/.profile <<'EOF'

export GOROOT=$HOME/go
export PATH="$PATH:$GOROOT/bin"

EOF


  echo "- Installing as user - Go"
  sudo -i -u pi /bin/sh -c "curl -S https://storage.googleapis.com/golang/go1.7.linux-armv6l.tar.gz | tar xz ~"
fi


echo "- Installing as user - setting up GOPATH"
cat >> /home/pi/.profile <<'EOF'

export GOPATH=$HOME
export PATH="$PATH:$GOPATH/bin"

EOF


echo "- Installing as user - dlibox"
sudo -i -u pi /bin/sh -c "go get -v github.com/maruel/dlibox/go/cmd/dlibox"


# At this point, defer control to the script in the repository.
/home/pi/src/github.com/maruel/dlibox/go/setup/host/finish_install.sh
