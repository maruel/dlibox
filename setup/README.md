# Setup

## Software

- Raspbian Jessie Lite
- Enable the SPI port on the Raspberry Pi via `raspi-config`.
- `sudo apt-get install libcap2-bin ntpdate`


## Auto-start on boot and auto-restart on scp

    sudo cp dotstar.service/etc/systemd/system
    sudo systemctl enable dotstar
    sudo service dotstar start
