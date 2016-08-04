#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Find the raspberry pis on the local network.

avahi-browse -t _workstation._tcp | grep IPv4 | cut -f 4 -d ' ' | sort | uniq
