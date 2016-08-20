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
Doc](https://godoc.org/github.com/maruel/dlibox/go/rpi?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/rpi).


### Functional

Print the pins per special functionality:

    $ pins -f
    I2C_SCL1 : GPIO3
    I2C_SDA1 : GPIO2
    IR_IN    : GPIO6
    IR_OUT   : GPIO5
    PWM0_OUT : GPIO40
    PWM1_OUT : GPIO41
    SPI0_CLK : GPIO11
    SPI0_MISO: GPIO9
    SPI0_MOSI: GPIO10
    SPI1_CLK : GPIO21
    SPI1_MISO: GPIO19
    SPI1_MOSI: GPIO20
    UART_RXD0: GPIO15
    UART_TXD0: GPIO14


### GPIO

Print the pins per GPIO number:

    $ ./pins -g
    GPIO0 : In/High
    GPIO1 : In/High
    GPIO2 : I2C_SDA1
    GPIO3 : I2C_SCL1
    GPIO4 : In/High
    GPIO5 : IR_OUT
    GPIO6 : IR_IN
    GPIO7 : Out/High
    GPIO8 : Out/High
    GPIO9 : SPI0_MISO
    GPIO10: SPI0_MOSI
    GPIO11: SPI0_CLK
    GPIO12: In/Low
    GPIO13: In/Low
    GPIO14: UART_TXD0
    GPIO15: UART_RXD0
    GPIO16: In/High
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
internal lookup table then query each pin.

    $ ./pins -h
    Header    Func  Name  Pos Pos  Name   Func
    P1:             V3_3  1 x x 2  V5     
         I2C_SDA1  GPIO2  3 x x 4  V5     
         I2C_SCL1  GPIO3  5 x x 6  GROUND 
          In/High  GPIO4  7 x x 8  GPIO14 UART_TXD0
                  GROUND  9 x x 10 GPIO15 UART_RXD0
           In/Low GPIO17 11 x x 12 GPIO18 Out/High
           In/Low GPIO27 13 x x 14 GROUND 
           In/Low GPIO22 15 x x 16 GPIO23 In/Low
                    V3_3 17 x x 18 GPIO24 In/Low
        SPI0_MOSI GPIO10 19 x x 20 GROUND 
        SPI0_MISO  GPIO9 21 x x 22 GPIO25 In/Low
         SPI0_CLK GPIO11 23 x x 24 GPIO8  Out/High
                  GROUND 25 x x 26 GPIO7  Out/High
          In/High  GPIO0 27 x x 28 GPIO1  In/High
           IR_OUT  GPIO5 29 x x 30 GROUND 
            IR_IN  GPIO6 31 x x 32 GPIO12 In/Low
           In/Low GPIO13 33 x x 34 GROUND 
        SPI1_MISO GPIO19 35 x x 36 GPIO16 In/High
           In/Low GPIO26 37 x x 38 GPIO20 SPI1_MOSI
                  GROUND 39 x x 40 GPIO21 SPI1_CLK
    
    P5:           INVALID 1 x x 2 INVALID 
                  INVALID 3 x x 4 INVALID 
                  INVALID 5 x x 6 INVALID 
                  INVALID 7 x x 8 INVALID 
    
    AUDIO_LEFT  : GPIO41
    AUDIO_RIGHT : GPIO40
    HDMI_HOTPLUG: GPIO46


### Info

Queries the Raspberry Pi version:

    $ ./pins -i
    Version: 3  MaxSpeed: 1200Mhz
