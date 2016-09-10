#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Make sure the script is running as root and on Raspbian.

set -eu

# When run via /etc/rc.local, USER is not defined.
if [ "${USER:=root}" != "root" ]; then
  echo "This script must be run as root."
  exit 1
fi


# The idea is that this command will fail if not running on Raspbian, as a
# safety measure.
echo "- Testing if running on Raspbian"
grep raspbian /etc/os-release > /dev/null
