# pins

Prints bcm238x pins.

## Examples

* Use `pins --help` for help
* Use `-n` to print pins that are not connected or in INVALID state
* Use `-a` to print everything at once


### Functional

    $ pins -f
    GPCLK1   : GPIO42
    GPCLK2   : GPIO43
    I2C_SCL1 : GPIO3
    I2C_SDA1 : GPIO2
    IR_IN    : GPIO18
    IR_OUT   : GPIO17
    PWM0_OUT : GPIO40
    PWM1_OUT : GPIO41
    SPI0_CLK : GPIO11
    SPI0_MISO: GPIO9
    SPI0_MOSI: GPIO10
    UART_RXD0: GPIO33
    UART_TXD0: GPIO32


### GPIO

    $ ./pins -g
    GPIO0 : In   High
    GPIO1 : In   High
    GPIO2 : Alt0
    GPIO3 : Alt0
    GPIO4 : In   High
    GPIO5 : In   High
    GPIO6 : In   High
    GPIO7 : Out  High
    GPIO8 : Out  High
    GPIO9 : Alt0
    GPIO10: Alt0
    GPIO11: Alt0
    GPIO12: In   Low
    GPIO13: In   Low
    GPIO14: In   Low
    GPIO15: In   High
    GPIO16: In   Low
    GPIO17: Out  Low
    GPIO18: In   High
    GPIO19: In   Low
    GPIO20: In   Low
    GPIO21: In   Low
    GPIO22: In   Low
    GPIO23: In   Low
    GPIO24: In   Low
    GPIO25: In   Low
    GPIO26: In   Low
    GPIO27: In   Low
    GPIO40: Alt0
    GPIO45: In   Low
    GPIO46: In   Low


### Hardware

This uses an internal lookup table.

    $ ./pins -h
    P1_1        : V3_3
    P1_2        : V5
    P1_3        : GPIO2
    P1_4        : V5
    P1_5        : GPIO3
    P1_6        : GROUND
    P1_7        : GPIO4
    P1_8        : GPIO14
    P1_9        : GROUND
    P1_10       : GPIO15
    P1_11       : GPIO17
    P1_12       : GPIO18
    P1_13       : GPIO27
    P1_14       : GROUND
    P1_15       : GPIO22
    P1_16       : GPIO23
    P1_17       : V3_3
    P1_18       : GPIO24
    P1_19       : GPIO10
    P1_20       : GROUND
    P1_21       : GPIO9
    P1_22       : GPIO25
    P1_23       : GPIO11
    P1_24       : GPIO8
    P1_25       : GROUND
    P1_26       : GPIO7
    P1_27       : GPIO0
    P1_28       : GPIO1
    P1_29       : GPIO5
    P1_30       : GROUND
    P1_31       : GPIO6
    P1_32       : GPIO12
    P1_33       : GPIO13
    P1_34       : GROUND
    P1_35       : GPIO19
    P1_36       : GPIO16
    P1_37       : GPIO26
    P1_38       : GPIO20
    P1_39       : GROUND
    P1_40       : GPIO21
    AUDIO_LEFT  : GPIO45
    AUDIO_RIGHT : GPIO40
    HDMI_HOTPLUG: GPIO46
