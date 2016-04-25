# Setup

## Software

  - [Raspbian Jessie Lite](https://www.raspberrypi.org/downloads/raspbian/).
  - Enable the SPI port on the Raspberry Pi via `raspi-config`.
  - Download and install [go1.6.2.linux-armv6l.tar.gz](https://golang.org/dl/) or the latest version.
  - Run a few commands


        sudo apt-get install libcap2-bin ntpdate
        export GOPATH=$HOME
        export PATH="$GOPATH:$PATH"
        echo 'export GOPATH=$HOME' >> $HOME/.bash_aliases
        echo 'export PATH="$GOPATH:$PATH"' >> $HOME/.bash_aliases
        go get github.com/maruel/dotstar/cmd/dotstar


## Auto-start on boot and auto-restart on scp

    sudo cp dotstar.service /etc/systemd/system
    sudo systemctl enable dotstar
    sudo service dotstar start

Then you can scp a new version in, mv over and the executable will be restarted
automatically.
