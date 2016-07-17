# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Run "source env.sh" to set the environment variables.

if [[ $_ == $0 ]]; then
  echo "Please source the script instead of running it:"
  echo "  source env.sh"
  exit 1
fi

# Change it to your liking.
SMING_BASE="$HOME/src"

export ESP_HOME="$SMING_BASE/esp-open-sdk"
export SMING_HOME="$SMING_BASE/Sming/Sming"

if [ ! -d "$ESP_HOME" ]; then
  echo "Please follow instructions to create $ESP_HOME"
elif [ ! -d "$SMING_HOME" ]; then
  echo "Please follow instructions to create $SMING_HOME"
else
  export PATH="$PATH:$SMING_BASE/esptool2"
fi
