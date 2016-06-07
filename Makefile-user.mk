# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Default: 115200
# COM_SPEED = 921600
# Running at 921.6kbps makes flashing flake out occasionally but it's
# significantly faster.
COM_SPEED_ESPTOOL = 921600

RBOOT_ENABLED ?= 1

SPI_SIZE ?= 4M

# TODO(maruel): Defaults to 40. Why not 80Mhz?
#SPI_SPEED = 80
