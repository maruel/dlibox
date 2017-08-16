# Required Hardware


## Computing

A few of:

- [C.H.I.P.](https://getchip.com); it has its own integrated LiPo battery
  charger, which is useful to keep as the controller that stays alive even in
  loss of power, assuming you also have backup power on your wifi router.
- Raspberry Pi 1/2/3 or Zero Wireless. The main advantage of the RPi3 is that it
  boots faster.
- Raspberry Pi case to stick a breadboard over.


## LEDs

- I do recommend APA-102 (also called Dotstar) from [iPixel LED Light
  Co](http://www.ipixelleds.com) for quality strips.
  - One of the nice thing from these is that the chip is rated for -40°C so
    buying a weatherproof strip + long cables is worth it.
- One 74AHCT125N to safely upgrade the SPI signal from 3.3V to 5V.
- One large (>=10A) 5V power supply per ~400 LEDs.
- Power wires for the same length of the strip to lower resistance (I used 18
  AWG).
- Small breadboard that can stick on the Raspberry Pi case.
- Small wires to connect from the Pi to the breadboard.
