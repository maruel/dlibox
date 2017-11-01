#!/bin/sh
set -eu
mosquitto_sub -t "#" -v -h dlibox
