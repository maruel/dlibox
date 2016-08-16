// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package rpi contains Raspberry Pi and BCM238x interfacing code.
//
// Requires Raspbian Jessie.
//
// Contains both low level GPIO logic (including edge triggering) and higher
// level communication like InfraRed, SPI (2 buses) and I²C (2 buses).
//
// No code in this module is "thread-safe".
//
// Configuration
//
// The pins function can be affected by device overlays as defined in
// /boot/config.txt. The full documentation of overlays is at
// https://github.com/raspberrypi/firmware/blob/master/boot/overlays/README.
// Documentation for the file format at
// https://www.raspberrypi.org/documentation/configuration/device-tree.md#part3
//
// I²C
//
// The BCM238x has 2 I²C buses.
//
// - /dev/i2c-1 can be enabled with:
//    dtparam=i2c=on
// - /dev/i2c-0 can be enabled with the following but be warned that it
// conflicts with HAT EEPROM detection at boot
// https://github.com/raspberrypi/hats
//    dtparam=i2c_vc=on
//    dtoverlay=i2c0-bcm2708 (Confirm?)
//
// I2S
//
// Can be enabled with:
//     dtparam=i2s=on
//
// IR
//
// Exposed as /dev/lirc0. Can be enabled with:
//     dtoverlay=lirc-rpi,gpio_out_pin=17,gpio_in_pin=18,gpio_in_pull=down
//
// Default pins 17 and 18 clashes with SPI1 so change the pin if you plan to
// enable both SPI buses.
//
// IR/Debugging
//
//     # Detect your remote
//     irrecord -a -d /var/run/lirc/lircd ~/lircd.conf
//     # Grep for key names you found to find the remote in the remotes library
//     grep -R '<hex value>' /usr/share/lirc/remotes/
//     # Listen and send command to the server
//     nc -U /var/run/lirc/lircd
//     # List all valid key names
//     irrecord -l
//     grep -hoER '(BTN|KEY)_\w+' /usr/share/lirc/remotes | sort | uniq | less
//
// Keys are listed at
// http://www.lirc.org/api-docs/html/input__map_8inc_source.html
//
// Source at:
// https://github.com/raspberrypi/linux/blob/rpi-4.8.y/drivers/staging/media/lirc/lirc_rpi.c
// Someone made a version that supports multiple devices:
// https://github.com/bengtmartensson/lirc_rpi
//
// PWM
//
// To take back control to use as general purpose PWM, comment out the
// following line:
//     dtparam=audio=on
//
// SPI
//
// The BCM238x has 3 SPI buses but only two are usable.
//
// - /dev/spidev0.0 and /dev/spidev0.1 can be enabled with:
//     dtparam=spi=on
// - /dev/spidev1.0 can be enabled with:
//     dtoverlay=spi1-1cs
// On rPi3, bluetooth must be disabled with:
//     dtoverlay=pi3-disable-bt
// and bluetooth UART service needs to be disabled with:
//     sudo systemctl disable hciuart
//
// UART
//
// Kernel boot messages go to the UART (0 or 1, depending on Pi version) at
// 115200 bauds.
//
// On Rasberry Pi 1 and 2, UART0 is used.
//
// On Raspberry Pi 3, UART0 is connected to bluetooth so the console is
// connected to UART1 instead. Disabling bluetooth also reverts to use UART0
// and not UART1.
//
// UART0 can be disabled with:
//     dtparam=uart0=off
// UART1 can be enabled with:
//     dtoverlay=uart1
package rpi
