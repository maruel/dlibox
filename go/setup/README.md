# Setup

- [prep/](prep) contains the scripts to run on a workstation to prepare an
  image.
- [host/](host) contains the scripts to run on the host, the Raspberry Pi.


## As a fresh image

### Raspberry Pi

The simplest is to create a fresh Raspbian Jessie Lite image using `flash.sh`.
You just put the card in your Raspberry Pi and it will initializes itself on
first boot.

- Make sure you have an ssh key configured.
- Insert a SDCard and note the device path, e.g. /dev/sdX where X is a letter.
- Run `flash.sh /dev/sdX <wifi ssid>`

This script generates an image that leaves 390Mb free on a 2Gb SDCard.


### Orange Pi / Pine64 / Banana Pi

Flash the image yourself with [Armbian](http://www.armbian.com/), chose Jessie
Server.

To initiate the bootstrapping process over ssh, run:

    ./armbian.sh <hostname>

Read the script for more details.


### SSH key

- Make sure your `.ssh/config` has the proper config to push to the account on
  which you want the service to run on. For example:

      Host dlibox-*
        StrictHostKeyChecking no
        User pi


## On existing installation

Configure your Raspberry Pi with everything necessary and start dlibox:

    make HOST=mypi setup

Push a new version:

    make HOST=mypi push

`HOST` defaults to `dlibox`.


### Manual

Read the scripts and do the same. In short they do:
- Install git and [Go](https://golang.org/dl/).
- `go get github.com/maruel/dlibox/go/pio/cmd/... github.com/maruel/dlibox/go/cmd/...`
- `sudo $GOPATH/src/github.com/maruel/dlibox/go/setup/host/install_systemd.sh`

Anytime you `go install github.com/maruel/dlibox/go/cmd/dlibox`, systemd will
restart dlibox automatically.


## Debugging


### mDNS

[Bonjour
Browser](https://play.google.com/store/apps/details?id=com.grokkt.android.bonjour)
is a nice Android app to debug mDNS broadcasts.


### MQTT

[MQTTLens](https://chrome.google.com/webstore/detail/mqttlens/hemojaaeigabkbcookmlgmdigohjobjm)
is a Google Chrome app to debug messages on a MQTT server.


### Logs

Look at the logs on the dlibox server:

    sudo journalctl -u dlibox
    # For streaming:
    sudo journalctl -f -u dlibox


### InfraRed receiver

`irw` will print decoded messages by [lirc](http://www.lirc.org/) via the
[Raspbian specific lirc-rpi kernel
module](https://github.com/raspberrypi/firmware/blob/master/boot/overlays/README).
See [rpi/ir.go](../rpi/ir.go) for more details.
