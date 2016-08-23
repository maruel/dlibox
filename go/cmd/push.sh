#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu
cd "$(dirname $0)"

if [ "$#" != 2 ]; then
  echo "usage: $0 <dlibox hostname> <tool name>"
  exit 1
fi

NAME="${2%/}"
HOST="$1"

cd "$NAME"
GOOS=linux GOARCH=arm go build .
scp "$NAME" "$HOST:."
rm "$NAME"
