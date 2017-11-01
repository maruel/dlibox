#!/bin/sh

set -eu

GOOS=linux GOARCH=arm go test -i ./cmd/dlibox
HOSTS="dlibox raspberrypi-e59e raspberrypi-73f5 raspberrypi-681e"
for i in $HOSTS; do
  echo $i
  #ssh $i "sudo systemctl disable dlibox_update"
  ./cmd/push.sh $i
done
