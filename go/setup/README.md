# Setup

The script `setup.sh` is meant to be used along
[github.com/periph/bootstrap](https://github.com/periph/bootstrap). This is
still a work in progress.


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
