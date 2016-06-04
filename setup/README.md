# Setup

## Software

- Install [Raspbian Jessie
  Lite](https://www.raspberrypi.org/downloads/raspbian/) and make sure you can
  ssh to the Raspberry Pi.
- Enable the SPI port on the Raspberry Pi via
  [`raspi-config`](https://www.raspberrypi.org/documentation/configuration/raspi-config.md).


### Automated

Configure your Raspberry Pi with everything necessary and start dlibox:

    make HOST=mypi setup

Push a new version:

    make HOST=mypi push

`HOST` defaults to `raspberrypi1`.


### Manual

This enables building dlibox from the rPi itself. It's a bit slow on a rPi 1
but it's totally acceptable on a rPi 2 or rPi 3.

_Note:_ Replace the URL below with the [latest version](https://golang.org/dl/).

    sudo apt-get install libcap2-bin ntpdate
    curl https://storage.googleapis.com/golang/go1.6.2.linux-armv6l.tar.gz | tar xz
    echo 'export GOPATH=$HOME' >> $HOME/.bash_aliases
    echo 'export GOROOT=$HOME/go' >> $HOME/.bash_aliases
    echo 'export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"' >> $HOME/.bash_aliases
    source $HOME/.bash_aliases
    go get github.com/maruel/dlibox-go/cmd/dlibox
    # If you plan to do edit-compile, you can precompile all dependencies:
    go test -i github.com/maruel/dlibox-go/cmd/dlibox

Set it up to auto-start on boot and auto-restart on scp:

    sudo cp $GOPATH/src/github.com/maruel/dliboxt/setup/dlibox.service /etc/systemd/system
    # Edit the file with the right user and path
    sudo vi $GOPATH/src/github.com/maruel/dliboxt/setup/dlibox.service
    sudo systemctl enable dlibox
    sudo service dlibox start

Anytime you `go install github.com/maruel/dlibox-go/cmd/dlibox`, systemd will
restart dlibox automatically.


## Logs

    sudo journalctl -u dlibox
