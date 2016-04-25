# Setup

## Software

- [Raspbian Jessie Lite](https://www.raspberrypi.org/downloads/raspbian/).
- Enable the SPI port on the Raspberry Pi via `raspi-config`.
- Run the following commands:

Note: Replace the URL below with the [latest version](https://golang.org/dl/).

    sudo apt-get install libcap2-bin ntpdate
    curl https://storage.googleapis.com/golang/go1.6.2.linux-armv6l.tar.gz | tar xz
    echo 'export GOPATH=$HOME' >> $HOME/.bash_aliases
    echo 'export GOROOT=$HOME/go' >> $HOME/.bash_aliases
    echo 'export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"' >> $HOME/.bash_aliases
    source $HOME/.bash_aliases
    go get github.com/maruel/dotstar/cmd/dotstar


## Auto-start on boot and auto-restart on scp

    sudo cp dotstar.service /etc/systemd/system
    sudo systemctl enable dotstar
    sudo service dotstar start

Then you can scp a new version in, mv over and the executable will be restarted
automatically.
