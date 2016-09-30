# pio-info

Prints the lists of drivers that were loaded, the ones skipped and the one that
failed to load, if any.

* Looking for the GPIO pins per functionality? Look at
  [gpio-list](../gpio-list).
* Looking for the location of the pin on the header to connect your GPIO? Look
  at [headers-list](../headers-list).


## Example

    $ pio-info
    Using drivers:
      - bcm283x
      - sysfs-leds
      - rpi
      - sysfs-gpio
      - sysfs-spi
      - sysfs-i2c
    Drivers skipped:
      - allwinner
      - allwinner_pl
      - pine64
