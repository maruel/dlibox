#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# For an already setup Armbian host, initializes it remotely.
# Do not fetch like for Raspbian since there are many images. Do everything over
# ssh.
#
# Assumptions:
# - ssh root@host once to setup password and created a user.
# - Because armbian doesn't have avahi preinstalled, it can'be be found with
#   ./find.sh. :(
# - Ethernet connection is used for the bootstrapping process. This means that
#   wifi is setup afterward; unlike Raspbian which has its image modified
#   upfront, so wired network is not needed during bootstrap.

set -eu
cd "$(dirname $0)"

if [ "$#" -ne 1 ]; then
  echo "Setups an armbian host via ssh."
  exit 1
fi


# TODO(maruel): There are 3 things that are done via ./steps/setup_firstboot.sh
# that needs to be done here:
# - Setup ssh key
# - Disable password based ssh authentication
# - Setup wifi
ssh root@$1 "mkdir .ssh"
scp ~/.ssh/authorized_keys root@$1:.ssh/

# TODO(maruel): Things that needs to be done due to differences between Raspbian
# and Armbian:
# - Disable ssh as root
# - Remove root password
# - Create 'pi' user
# - Make sudo passwordless (which is crazy, so we should probably change
#   Raspbian behavior instead).
scp support/dlibox_firstboot.sh root@$1:.
# TODO(maruel): Do not run apt-get update here since it was already done as part
# of the initial boot. Overall it's just ~30s saving.
ssh root@$1 "bash dlibox_firstboot.sh"
