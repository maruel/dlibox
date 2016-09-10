# host

This directory contains files meant to be run on a Raspberry Pi/Orange
Pi/Pine64/etc to configure it to run as a dlibox.

In an installation aborts midway, let's say due to network failure, you can
resume by manually executing the steps in `finish_install.sh`


## Display

If you use a 5" display and the right side is corrupted, mount the card on a
host and run [../prep/enable_5inch.sh](../prep/enable_5inch.sh) to enable the
full width.

If you get a blank screen, mount the card on a host and edit
`/etc/systemd/system/disable_hdmi.service` to replace

    ExecStart=/opt/vc/bin/tvservice -o

with

    ExecStart=/opt/vc/bin/tvservice -p


### Using X

The image used doesn't have X installed but if you use Jessie instead of Jessie
Lite, X won't start because `install_dependencies.sh` reduces the GPU reserved
memory down to 16Mb so X can't start. Undo with:

    raspi-config nonint do_memory_split 64
