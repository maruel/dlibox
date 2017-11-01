#!/bin/sh
set -eu
mosquitto_pub -h dlibox -t dlibox/raspberrypi-e59e/pir/pir -m foo
