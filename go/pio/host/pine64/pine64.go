// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pine64

import (
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/a64"
	"github.com/maruel/dlibox/go/pio/host/pins"
)

// Version is the board version. Only reports as 1 for now.
var Version int = 1

var (
	PI2_1  host.Pin = pins.V3_3   // 3.3 volt; max 40mA
	PI2_2  host.Pin = pins.V5     // 5 volt (before filtering)
	PI2_3  host.Pin = a64.PH3     //
	PI2_4  host.Pin = pins.V5     //
	PI2_5  host.Pin = a64.PH2     //
	PI2_6  host.Pin = pins.GROUND //
	PI2_7  host.Pin = a64.PL10    //
	PI2_8  host.Pin = a64.PB0     //
	PI2_9  host.Pin = pins.GROUND //
	PI2_10 host.Pin = a64.PB1     //
	PI2_11 host.Pin = a64.PC7     //
	PI2_12 host.Pin = a64.PC8     //
	PI2_13 host.Pin = a64.PH9     //
	PI2_14 host.Pin = pins.GROUND //
	PI2_15 host.Pin = a64.PC12    //
	PI2_16 host.Pin = a64.PC13    //
	PI2_17 host.Pin = pins.V3_3   //
	PI2_18 host.Pin = a64.PC14    //
	PI2_19 host.Pin = a64.PC0     //
	PI2_20 host.Pin = pins.GROUND //
	PI2_21 host.Pin = a64.PC1     //
	PI2_22 host.Pin = a64.PC15    //
	PI2_23 host.Pin = a64.PC2     //
	PI2_24 host.Pin = a64.PC3     //
	PI2_25 host.Pin = pins.GROUND //
	PI2_26 host.Pin = a64.PH7     //
	PI2_27 host.Pin = a64.PL9     //
	PI2_28 host.Pin = a64.PL8     //
	PI2_29 host.Pin = a64.PH5     //
	PI2_30 host.Pin = pins.GROUND //
	PI2_31 host.Pin = a64.PH6     //
	PI2_32 host.Pin = a64.PC4     //
	PI2_33 host.Pin = a64.PC5     //
	PI2_34 host.Pin = pins.GROUND //
	PI2_35 host.Pin = a64.PC9     //
	PI2_36 host.Pin = a64.PC6     //
	PI2_37 host.Pin = a64.PC16    //
	PI2_38 host.Pin = a64.PC10    //
	PI2_39 host.Pin = pins.GROUND //
	PI2_40 host.Pin = a64.PC11    //

	EULER_1  host.Pin = pins.V3_3        //
	EULER_2  host.Pin = pins.DC_IN       //
	EULER_3  host.Pin = pins.BAT_PLUS    //
	EULER_4  host.Pin = pins.DC_IN       //
	EULER_5  host.Pin = pins.TEMP_SENSOR //
	EULER_6  host.Pin = pins.GROUND      //
	EULER_7  host.Pin = pins.IR_RX       //
	EULER_8  host.Pin = pins.V5          //
	EULER_9  host.Pin = pins.GROUND      //
	EULER_10 host.Pin = a64.PH8          //
	EULER_11 host.Pin = a64.PB3          //
	EULER_12 host.Pin = a64.PB4          //
	EULER_13 host.Pin = a64.PB5          //
	EULER_14 host.Pin = pins.GROUND      //
	EULER_15 host.Pin = a64.PB6          //
	EULER_16 host.Pin = a64.PB7          //
	EULER_17 host.Pin = pins.V3_3        //
	EULER_18 host.Pin = a64.PD4          //
	EULER_19 host.Pin = a64.PD2          //
	EULER_20 host.Pin = pins.GROUND      //
	EULER_21 host.Pin = a64.PD3          //
	EULER_22 host.Pin = a64.PD5          //
	EULER_23 host.Pin = a64.PD1          //
	EULER_24 host.Pin = a64.PD0          //
	EULER_25 host.Pin = pins.GROUND      //
	EULER_26 host.Pin = a64.PD6          //
	EULER_27 host.Pin = a64.PB2          //
	EULER_28 host.Pin = a64.PD7          //
	EULER_29 host.Pin = a64.PB8          //
	EULER_30 host.Pin = a64.PB9          //
	EULER_31 host.Pin = pins.EAROUTP     //
	EULER_32 host.Pin = pins.EAROUT_N    //
	EULER_33 host.Pin = pins.INVALID     //
	EULER_34 host.Pin = pins.GROUND      //

	EXP_1  host.Pin = pins.V3_3        //
	EXP_2  host.Pin = a64.PL7          //
	EXP_3  host.Pin = pins.CHARGER_LED //
	EXP_4  host.Pin = pins.RESET       //
	EXP_5  host.Pin = pins.PWR_SWITCH  //
	EXP_6  host.Pin = pins.GROUND      //
	EXP_7  host.Pin = a64.PB8          //
	EXP_8  host.Pin = a64.PB9          //
	EXP_9  host.Pin = pins.GROUND      //
	EXP_10 host.Pin = pins.KEY_ADC     //

	WIFI_BT_1  host.Pin = pins.GROUND   //
	WIFI_BT_2  host.Pin = a64.PG6       //
	WIFI_BT_3  host.Pin = a64.PG0       //
	WIFI_BT_4  host.Pin = a64.PG7       //
	WIFI_BT_5  host.Pin = pins.GROUND   //
	WIFI_BT_6  host.Pin = a64.PG8       //
	WIFI_BT_7  host.Pin = a64.PG1       //
	WIFI_BT_8  host.Pin = a64.PG9       //
	WIFI_BT_9  host.Pin = a64.PG2       //
	WIFI_BT_10 host.Pin = a64.PG10      //
	WIFI_BT_11 host.Pin = a64.PG3       //
	WIFI_BT_12 host.Pin = a64.PG11      //
	WIFI_BT_13 host.Pin = a64.PG4       //
	WIFI_BT_14 host.Pin = a64.PG12      //
	WIFI_BT_15 host.Pin = a64.PG5       //
	WIFI_BT_16 host.Pin = a64.PG13      //
	WIFI_BT_17 host.Pin = a64.PL2       //
	WIFI_BT_18 host.Pin = pins.GROUND   //
	WIFI_BT_19 host.Pin = a64.PL3       //
	WIFI_BT_20 host.Pin = a64.PL5       //
	WIFI_BT_21 host.Pin = pins.X32KFOUT //
	WIFI_BT_22 host.Pin = a64.PL5       //
	WIFI_BT_23 host.Pin = pins.GROUND   //
	WIFI_BT_24 host.Pin = a64.PL6       //
	WIFI_BT_25 host.Pin = pins.VCC      //
	WIFI_BT_26 host.Pin = pins.IOVCC    //

	AUDIO_LEFT  host.Pin = pins.INVALID // TODO(maruel): Figure out, is that EAROUT?
	AUDIO_RIGHT host.Pin = pins.INVALID //
)
