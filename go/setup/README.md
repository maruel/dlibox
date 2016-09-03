# Setup

## As a fresh image

The simplest is to create a fresh Raspbian Jessie Lite image using `flash.sh`.
You just put the card in your Raspberry Pi and it will initializes itself on
first boot.

- Make sure you have an ssh key configured.
- Insert a SDCard and note the device path, e.g. /dev/sdX where X is a letter.
- Run `flash.sh /dev/sdX <wifi ssid>`

This script generates an image that leaves 390Mb free on a 2Gb SDCard.


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

This enables building dlibox from the rPi itself. It's a bit slow on a rPi 1
but it's totally acceptable on a rPi 2 or rPi 3.

_Note:_ Replace the URL below with the [latest version](https://golang.org/dl/).

    cd
    curl https://storage.googleapis.com/golang/go1.7.linux-armv6l.tar.gz | tar xz
    echo 'export GOPATH=$HOME' >> $HOME/.profile
    echo 'export GOROOT=$HOME/go' >> $HOME/.profile
    echo 'export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"' >> $HOME/.profile
    source $HOME/.profile
    go get github.com/maruel/dlibox/go/cmd/dlibox
    # If you plan to do edit-compile, you can precompile all dependencies:
    go test -i github.com/maruel/dlibox/go/cmd/dlibox
    # Run apt-get install ...
    sudo $GOPATH/src/github.com/maruel/dlibox/go/setup/support/install_dependencies.sh

Set it up to auto-start on boot and auto-restart on scp:

    sudo $GOPATH/src/github.com/maruel/dlibox/go/setup/support/install_systemd.sh
    sudo service dlibox start

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
