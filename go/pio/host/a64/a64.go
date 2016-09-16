// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package a64

var Pins []Pin

// Page 23~24
// Each pin supports 6 functions.

type Pin uint8

// http://forum.pine64.org/showthread.php?tid=474
// about number calculation.
const (
	PB0  Pin = 0x00 // Z UART2_TX, -, JTAG_MS0, -, PB_EINT0
	PB1  Pin = 0x01 // Z, UART2_RX, -, JTAG_CK0, SIM_PWREN, PB_EINT1
	PB2  Pin = 0x02 //
	PB3  Pin = 0x03 //
	PB4  Pin = 0x04 //
	PB5  Pin = 0x05 //
	PB6  Pin = 0x06 //
	PB7  Pin = 0x07 //
	PB8  Pin = 0x08 //
	PB9  Pin = 0x09 //
	PC0  Pin = 0x10 //
	PC1  Pin = 0x11 //
	PC2  Pin = 0x12 //
	PC3  Pin = 0x13 //
	PC4  Pin = 0x14 //
	PC5  Pin = 0x15 //
	PC6  Pin = 0x16 //
	PC7  Pin = 0x17 //
	PC8  Pin = 0x18 //
	PC9  Pin = 0x19 //
	PC10 Pin = 0x1A //
	PC11 Pin = 0x1B //
	PC12 Pin = 0x1C //
	PC13 Pin = 0x1D //
	PC14 Pin = 0x1E //
	PC15 Pin = 0x1F //
	PC16 Pin = 0x20 //
	PD0  Pin = 0x30 //
	PD1  Pin = 0x31 //
	PD2  Pin = 0x32 //
	PD3  Pin = 0x33 //
	PD4  Pin = 0x34 //
	PD5  Pin = 0x35 //
	PD6  Pin = 0x36 //
	PD7  Pin = 0x37 //
	PD8  Pin = 0x38 //
	PD9  Pin = 0x39 //
	PD10 Pin = 0x3A //
	PD11 Pin = 0x3B //
	PD12 Pin = 0x3C //
	PD13 Pin = 0x3D //
	PD14 Pin = 0x3E //
	PD15 Pin = 0x3F //
	PD16 Pin = 0x40 //
	PD17 Pin = 0x41 //
	PD18 Pin = 0x42 //
	PD19 Pin = 0x43 //
	PD20 Pin = 0x44 //
	PD21 Pin = 0x45 //
	PD22 Pin = 0x46 //
	PD23 Pin = 0x47 //
	PD24 Pin = 0x48 //
	PE0  Pin = 0x50 //
	PE1  Pin = 0x51 //
	PE2  Pin = 0x52 //
	PE3  Pin = 0x53 //
	PE4  Pin = 0x54 //
	PE5  Pin = 0x55 //
	PE6  Pin = 0x56 //
	PE7  Pin = 0x57 //
	PE8  Pin = 0x58 //
	PE9  Pin = 0x59 //
	PE10 Pin = 0x5A //
	PE11 Pin = 0x5B //
	PE12 Pin = 0x5C //
	PE13 Pin = 0x5D //
	PE14 Pin = 0x5E //
	PE15 Pin = 0x5F //
	PE16 Pin = 0x60 //
	PE17 Pin = 0x61 //
	PF0  Pin = 0x70 //
	PF1  Pin = 0x71 //
	PF2  Pin = 0x72 //
	PF3  Pin = 0x73 //
	PF4  Pin = 0x74 //
	PF5  Pin = 0x75 //
	PF6  Pin = 0x76 //
	PG0  Pin = 0x80 //
	PG1  Pin = 0x81 //
	PG2  Pin = 0x82 //
	PG3  Pin = 0x83 //
	PG4  Pin = 0x84 //
	PG5  Pin = 0x85 //
	PG6  Pin = 0x86 //
	PG7  Pin = 0x87 //
	PG8  Pin = 0x88 //
	PG9  Pin = 0x89 //
	PG10 Pin = 0x8A //
	PG11 Pin = 0x8B //
	PG12 Pin = 0x8C //
	PG13 Pin = 0x8D //
	PH0  Pin = 0x90 //
	PH1  Pin = 0x91 //
	PH2  Pin = 0x92 //
	PH3  Pin = 0x93 //
	PH4  Pin = 0x94 //
	PH5  Pin = 0x95 //
	PH6  Pin = 0x96 //
	PH7  Pin = 0x97 //
	PH8  Pin = 0x98 //
	PH9  Pin = 0x99 //
	PH10 Pin = 0x9A //
	PH11 Pin = 0x9B //
	PL1  Pin = 0xA1 //
	PL2  Pin = 0xA2 //
	PL3  Pin = 0xA3 //
	PL4  Pin = 0xA4 //
	PL5  Pin = 0xA5 //
	PL6  Pin = 0xA6 //
	PL7  Pin = 0xA7 //
	PL8  Pin = 0xA8 //
	PL9  Pin = 0xA9 //
	PL10 Pin = 0xAA //
	PL11 Pin = 0xAB //
	PL12 Pin = 0xAC //
)

func (p Pin) String() string {
	// TODO(maruel): Add.
	return ""
}

func (p Pin) Number() int {
	return int(p)
}

func (p Pin) Function() string {
	// TODO(maruel): Add.
	return ""
}
