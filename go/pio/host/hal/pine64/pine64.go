// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pine64

import (
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/drivers/allwinner"
	"github.com/maruel/dlibox/go/pio/host/hal/pins"
)

// Version is the board version. Only reports as 1 for now.
var Version int = 1

var (
	P1_1  host.PinIO = pins.V3_3      // 3.3 volt; max 40mA
	P1_2  host.PinIO = pins.V5        // 5 volt (before filtering)
	P1_3  host.PinIO = allwinner.PH3  //
	P1_4  host.PinIO = pins.V5        //
	P1_5  host.PinIO = allwinner.PH2  //
	P1_6  host.PinIO = pins.GROUND    //
	P1_7  host.PinIO = allwinner.PL10 //
	P1_8  host.PinIO = allwinner.PB0  //
	P1_9  host.PinIO = pins.GROUND    //
	P1_10 host.PinIO = allwinner.PB1  //
	P1_11 host.PinIO = allwinner.PC7  //
	P1_12 host.PinIO = allwinner.PC8  //
	P1_13 host.PinIO = allwinner.PH9  //
	P1_14 host.PinIO = pins.GROUND    //
	P1_15 host.PinIO = allwinner.PC12 //
	P1_16 host.PinIO = allwinner.PC13 //
	P1_17 host.PinIO = pins.V3_3      //
	P1_18 host.PinIO = allwinner.PC14 //
	P1_19 host.PinIO = allwinner.PC0  //
	P1_20 host.PinIO = pins.GROUND    //
	P1_21 host.PinIO = allwinner.PC1  //
	P1_22 host.PinIO = allwinner.PC15 //
	P1_23 host.PinIO = allwinner.PC2  //
	P1_24 host.PinIO = allwinner.PC3  //
	P1_25 host.PinIO = pins.GROUND    //
	P1_26 host.PinIO = allwinner.PH7  //
	P1_27 host.PinIO = allwinner.PL9  //
	P1_28 host.PinIO = allwinner.PL8  //
	P1_29 host.PinIO = allwinner.PH5  //
	P1_30 host.PinIO = pins.GROUND    //
	P1_31 host.PinIO = allwinner.PH6  //
	P1_32 host.PinIO = allwinner.PC4  //
	P1_33 host.PinIO = allwinner.PC5  //
	P1_34 host.PinIO = pins.GROUND    //
	P1_35 host.PinIO = allwinner.PC9  //
	P1_36 host.PinIO = allwinner.PC6  //
	P1_37 host.PinIO = allwinner.PC16 //
	P1_38 host.PinIO = allwinner.PC10 //
	P1_39 host.PinIO = pins.GROUND    //
	P1_40 host.PinIO = allwinner.PC11 //

	EULER_1  host.PinIO = pins.V3_3        //
	EULER_2  host.PinIO = pins.DC_IN       //
	EULER_3  host.PinIO = pins.BAT_PLUS    //
	EULER_4  host.PinIO = pins.DC_IN       //
	EULER_5  host.PinIO = pins.TEMP_SENSOR //
	EULER_6  host.PinIO = pins.GROUND      //
	EULER_7  host.PinIO = pins.IR_RX       //
	EULER_8  host.PinIO = pins.V5          //
	EULER_9  host.PinIO = pins.GROUND      //
	EULER_10 host.PinIO = allwinner.PH8    //
	EULER_11 host.PinIO = allwinner.PB3    //
	EULER_12 host.PinIO = allwinner.PB4    //
	EULER_13 host.PinIO = allwinner.PB5    //
	EULER_14 host.PinIO = pins.GROUND      //
	EULER_15 host.PinIO = allwinner.PB6    //
	EULER_16 host.PinIO = allwinner.PB7    //
	EULER_17 host.PinIO = pins.V3_3        //
	EULER_18 host.PinIO = allwinner.PD4    //
	EULER_19 host.PinIO = allwinner.PD2    //
	EULER_20 host.PinIO = pins.GROUND      //
	EULER_21 host.PinIO = allwinner.PD3    //
	EULER_22 host.PinIO = allwinner.PD5    //
	EULER_23 host.PinIO = allwinner.PD1    //
	EULER_24 host.PinIO = allwinner.PD0    //
	EULER_25 host.PinIO = pins.GROUND      //
	EULER_26 host.PinIO = allwinner.PD6    //
	EULER_27 host.PinIO = allwinner.PB2    //
	EULER_28 host.PinIO = allwinner.PD7    //
	EULER_29 host.PinIO = allwinner.PB8    //
	EULER_30 host.PinIO = allwinner.PB9    //
	EULER_31 host.PinIO = pins.EAROUTP     //
	EULER_32 host.PinIO = pins.EAROUT_N    //
	EULER_33 host.PinIO = host.INVALID     //
	EULER_34 host.PinIO = pins.GROUND      //

	EXP_1  host.PinIO = pins.V3_3        //
	EXP_2  host.PinIO = allwinner.PL7    //
	EXP_3  host.PinIO = pins.CHARGER_LED //
	EXP_4  host.PinIO = pins.RESET       //
	EXP_5  host.PinIO = pins.PWR_SWITCH  //
	EXP_6  host.PinIO = pins.GROUND      //
	EXP_7  host.PinIO = allwinner.PB8    //
	EXP_8  host.PinIO = allwinner.PB9    //
	EXP_9  host.PinIO = pins.GROUND      //
	EXP_10 host.PinIO = pins.KEY_ADC     //

	WIFI_BT_1  host.PinIO = pins.GROUND    //
	WIFI_BT_2  host.PinIO = allwinner.PG6  //
	WIFI_BT_3  host.PinIO = allwinner.PG0  //
	WIFI_BT_4  host.PinIO = allwinner.PG7  //
	WIFI_BT_5  host.PinIO = pins.GROUND    //
	WIFI_BT_6  host.PinIO = allwinner.PG8  //
	WIFI_BT_7  host.PinIO = allwinner.PG1  //
	WIFI_BT_8  host.PinIO = allwinner.PG9  //
	WIFI_BT_9  host.PinIO = allwinner.PG2  //
	WIFI_BT_10 host.PinIO = allwinner.PG10 //
	WIFI_BT_11 host.PinIO = allwinner.PG3  //
	WIFI_BT_12 host.PinIO = allwinner.PG11 //
	WIFI_BT_13 host.PinIO = allwinner.PG4  //
	WIFI_BT_14 host.PinIO = allwinner.PG12 //
	WIFI_BT_15 host.PinIO = allwinner.PG5  //
	WIFI_BT_16 host.PinIO = allwinner.PG13 //
	WIFI_BT_17 host.PinIO = allwinner.PL2  //
	WIFI_BT_18 host.PinIO = pins.GROUND    //
	WIFI_BT_19 host.PinIO = allwinner.PL3  //
	WIFI_BT_20 host.PinIO = allwinner.PL5  //
	WIFI_BT_21 host.PinIO = pins.X32KFOUT  //
	WIFI_BT_22 host.PinIO = allwinner.PL5  //
	WIFI_BT_23 host.PinIO = pins.GROUND    //
	WIFI_BT_24 host.PinIO = allwinner.PL6  //
	WIFI_BT_25 host.PinIO = pins.VCC       //
	WIFI_BT_26 host.PinIO = pins.IOVCC     //

	AUDIO_LEFT  host.PinIO = host.INVALID // TODO(maruel): Figure out, is that EAROUT?
	AUDIO_RIGHT host.PinIO = host.INVALID //
)

// See headers.Headers for more info.
var Headers = map[string][][]host.PinIO{
	"P1": {
		{P1_1, P1_2},
		{P1_3, P1_4},
		{P1_5, P1_6},
		{P1_7, P1_8},
		{P1_9, P1_10},
		{P1_11, P1_12},
		{P1_13, P1_14},
		{P1_15, P1_16},
		{P1_17, P1_18},
		{P1_19, P1_20},
		{P1_21, P1_22},
		{P1_23, P1_24},
		{P1_25, P1_26},
		{P1_27, P1_28},
		{P1_29, P1_30},
		{P1_31, P1_32},
		{P1_33, P1_34},
		{P1_35, P1_36},
		{P1_37, P1_38},
		{P1_39, P1_20},
	},
	"EULER": {
		{EULER_1, EULER_2},
		{EULER_3, EULER_4},
		{EULER_5, EULER_6},
		{EULER_7, EULER_8},
		{EULER_9, EULER_10},
		{EULER_11, EULER_12},
		{EULER_13, EULER_14},
		{EULER_15, EULER_16},
		{EULER_17, EULER_18},
		{EULER_19, EULER_20},
		{EULER_21, EULER_22},
		{EULER_23, EULER_24},
		{EULER_25, EULER_26},
		{EULER_27, EULER_28},
		{EULER_29, EULER_30},
		{EULER_31, EULER_32},
		{EULER_33, EULER_34},
	},
	"EXP": {
		{EXP_1, EXP_2},
		{EXP_3, EXP_4},
		{EXP_5, EXP_6},
		{EXP_7, EXP_8},
		{EXP_9, EXP_10},
	},
	"WIFI_BT": {
		{WIFI_BT_1, WIFI_BT_2},
		{WIFI_BT_3, WIFI_BT_4},
		{WIFI_BT_5, WIFI_BT_6},
		{WIFI_BT_7, WIFI_BT_8},
		{WIFI_BT_9, WIFI_BT_10},
		{WIFI_BT_11, WIFI_BT_12},
		{WIFI_BT_13, WIFI_BT_14},
		{WIFI_BT_15, WIFI_BT_16},
		{WIFI_BT_17, WIFI_BT_18},
		{WIFI_BT_19, WIFI_BT_20},
		{WIFI_BT_21, WIFI_BT_22},
		{WIFI_BT_23, WIFI_BT_24},
		{WIFI_BT_25, WIFI_BT_26},
	},
	"AUDIO": {
		{AUDIO_LEFT},
		{AUDIO_RIGHT},
	},
}
