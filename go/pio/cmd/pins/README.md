# pins

Prints bcm238x pins.

## Examples

* Use `pins --help` for help
* Use `-n` to print pins that are not connected or in INVALID state
* Use `-a` to print everything at once

The followings were captured on a Raspberry Pi 3 with I2C1, SPI0 and SPI1
enabled, lirc (IR) enabled and Bluetooth disabled with the following in
`/boot/config.txt`:

    dtparam=i2c_arm=on
    dtparam=spi=on
    dtoverlay=lirc-rpi,gpio_out_pin=5,gpio_in_pin=6,gpio_in_pull=high
    dtoverlay=spi1-1cs
    dtoverlay=pi3-disable-bt

then running:

    sudo systemctl disable hciuart

For more information for enabling functional pins, see
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio/host/rpi?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio/host/rpi).


### Functional

Print the pins per special functionality:

    $ pins -f
    GPCLK1   : GPIO42
    GPCLK2   : GPIO43
    I2C1_SCL : GPIO3
    I2C1_SDA : GPIO2
    PWM0_OUT : GPIO40
    PWM1_OUT : GPIO41
    SPI0_CLK : GPIO11
    SPI0_MISO: GPIO9
    SPI0_MOSI: GPIO10
    SPI1_CLK : GPIO21
    SPI1_MISO: GPIO19
    SPI1_MOSI: GPIO20
    UART0_RXD: GPIO15
    UART0_TXD: GPIO14


### GPIO

Print the pins per GPIO number:

    $ ./pins -g
    GPIO0 : In/High
    GPIO1 : In/High
    GPIO2 : I2C1_SDA
    GPIO3 : I2C1_SCL
    GPIO4 : In/Low
    GPIO5 : Out/Low
    GPIO6 : In/High
    GPIO7 : Out/High
    GPIO8 : Out/Low
    GPIO9 : SPI0_MISO
    GPIO10: SPI0_MOSI
    GPIO11: SPI0_CLK
    GPIO12: In/High
    GPIO13: In/High
    GPIO14: UART0_TXD
    GPIO15: UART0_RXD
    GPIO16: In/Low
    GPIO17: In/Low
    GPIO18: Out/High
    GPIO19: SPI1_MISO
    GPIO20: SPI1_MOSI
    GPIO21: SPI1_CLK
    GPIO22: In/Low
    GPIO23: In/Low
    GPIO24: In/Low
    GPIO25: In/Low
    GPIO26: In/Low
    GPIO27: In/Low
    GPIO40: PWM0_OUT
    GPIO41: PWM1_OUT
    GPIO46: In/High


### Hardware

Print the pins per their hardware location on the headers. This uses an
internal lookup table then query each pin. Here's an example on a host with two
SPI host and lirc enabled:

    $ ./pins -h
    AUDIO: 2 pins
      Pos  Name    Func
      1    GPIO41  PWM1_OUT
      2    GPIO40  PWM0_OUT

    HDMI: 1 pins
      Pos  Name    Func
      1    GPIO46  In/High

    P1: 40 pins
           Func    Name  Pos  Pos  Name   Func
                   V3_3    1  2    V5    g
       I2C1_SDA   GPIO2    3  4    V5    g
       I2C1_SCL   GPIO3    5  6    GROUNDg
         In/Low   GPIO4    7  8    GPIO14 UART0_TXD
                 GROUND    9  10   GPIO15 UART0_RXD
         In/Low  GPIO17   11  12   GPIO18 Out/High
         In/Low  GPIO27   13  14   GROUNDg
         In/Low  GPIO22   15  16   GPIO23 In/Low
                   V3_3   17  18   GPIO24 In/Low
      SPI0_MOSI  GPIO10   19  20   GROUNDg
      SPI0_MISO   GPIO9   21  22   GPIO25 In/Low
       SPI0_CLK  GPIO11   23  24   GPIO8  Out/Low
                 GROUND   25  26   GPIO7  Out/High
        In/High   GPIO0   27  28   GPIO1  In/High
        Out/Low   GPIO5   29  30   GROUNDg
        In/High   GPIO6   31  32   GPIO12 In/High
        In/High  GPIO13   33  34   GROUNDg
      SPI1_MISO  GPIO19   35  36   GPIO16 In/Low
         In/Low  GPIO26   37  38   GPIO20 SPI1_MOSI
                 GROUND   39  40   GPIO21 SPI1_CLK

### Info

Queries the processor's maximum speed:

    $ ./pins -i
    MaxSpeed: 1200Mhz
