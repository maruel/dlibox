// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pine64

import "github.com/maruel/dlibox/go/pio/host/a64"

// Version is the board version. Only reports as 1 for now.
var Version int = 1

var (
	PI2_1  = a64.V3_3   // 3.3 volt; max 40mA
	PI2_2  = a64.V5     // 5 volt (before filtering)
	PI2_3  = a64.PH3    //
	PI2_4  = a64.V5     //
	PI2_5  = a64.PH2    //
	PI2_6  = a64.GROUND //
	PI2_7  = a64.PL10   //
	PI2_8  = a64.PB0    //
	PI2_9  = a64.GROUND //
	PI2_10 = a64.PB1    //
	PI2_11 = a64.PC7    //
	PI2_12 = a64.PC8    //
	PI2_13 = a64.PH9    //
	PI2_14 = a64.GROUND //
	PI2_15 = a64.PC12   //
	PI2_16 = a64.PC13   //
	PI2_17 = a64.V3_3   //
	PI2_18 = a64.PC14   //
	PI2_19 = a64.PC0    //
	PI2_20 = a64.GROUND //
	PI2_21 = a64.PC1    //
	PI2_22 = a64.PC15   //
	PI2_23 = a64.PC2    //
	PI2_24 = a64.PC3    //
	PI2_25 = a64.GROUND //
	PI2_26 = a64.PH7    //
	PI2_27 = a64.PL9    //
	PI2_28 = a64.PL8    //
	PI2_29 = a64.PH5    //
	PI2_30 = a64.GROUND //
	PI2_31 = a64.PH6    //
	PI2_32 = a64.PC4    //
	PI2_33 = a64.PC5    //
	PI2_34 = a64.GROUND //
	PI2_35 = a64.PC9    //
	PI2_36 = a64.PC6    //
	PI2_37 = a64.PC16   //
	PI2_38 = a64.PC10   //
	PI2_39 = a64.GROUND //
	PI2_40 = a64.PC11   //

	EULER_1  = a64.V3_3        //
	EULER_2  = a64.DC_IN       //
	EULER_3  = a64.BAT_PLUS    //
	EULER_4  = a64.DC_IN       //
	EULER_5  = a64.TEMP_SENSOR //
	EULER_6  = a64.GROUND      //
	EULER_7  = a64.IR_RX       //
	EULER_8  = a64.V5          //
	EULER_9  = a64.GROUND      //
	EULER_10 = a64.PH8         //
	EULER_11 = a64.PB3         //
	EULER_12 = a64.PB4         //
	EULER_13 = a64.PB5         //
	EULER_14 = a64.GROUND      //
	EULER_15 = a64.PB6         //
	EULER_16 = a64.PB7         //
	EULER_17 = a64.V3_3        //
	EULER_18 = a64.PD4         //
	EULER_19 = a64.PD2         //
	EULER_20 = a64.GROUND      //
	EULER_21 = a64.PD3         //
	EULER_22 = a64.PD5         //
	EULER_23 = a64.PD1         //
	EULER_24 = a64.PD0         //
	EULER_25 = a64.GROUND      //
	EULER_26 = a64.PD6         //
	EULER_27 = a64.PB2         //
	EULER_28 = a64.PD7         //
	EULER_29 = a64.PB8         //
	EULER_30 = a64.PB9         //
	EULER_31 = a64.EAROUTP     //
	EULER_32 = a64.EAROUT_N    //
	EULER_33 = a64.INVALID     //
	EULER_34 = a64.GROUND      //

	EXP_1  = a64.V3_3        //
	EXP_2  = a64.PL7         //
	EXP_3  = a64.CHARGER_LED //
	EXP_4  = a64.RESET       //
	EXP_5  = a64.PWR_SWITCH  //
	EXP_6  = a64.GROUND      //
	EXP_7  = a64.PB8         //
	EXP_8  = a64.PB9         //
	EXP_9  = a64.GROUND      //
	EXP_10 = a64.KEY_ADC     //

	WIFI_BT_1  = a64.GROUND   //
	WIFI_BT_2  = a64.PG6      //
	WIFI_BT_3  = a64.PG0      //
	WIFI_BT_4  = a64.PG7      //
	WIFI_BT_5  = a64.GROUND   //
	WIFI_BT_6  = a64.PG8      //
	WIFI_BT_7  = a64.PG1      //
	WIFI_BT_8  = a64.PG9      //
	WIFI_BT_9  = a64.PG2      //
	WIFI_BT_10 = a64.PG10     //
	WIFI_BT_11 = a64.PG3      //
	WIFI_BT_12 = a64.PG11     //
	WIFI_BT_13 = a64.PG4      //
	WIFI_BT_14 = a64.PG12     //
	WIFI_BT_15 = a64.PG5      //
	WIFI_BT_16 = a64.PG13     //
	WIFI_BT_17 = a64.PL2      //
	WIFI_BT_18 = a64.GROUND   //
	WIFI_BT_19 = a64.PL3      //
	WIFI_BT_20 = a64.PL5      //
	WIFI_BT_21 = a64.X32KFOUT //
	WIFI_BT_22 = a64.PL5      //
	WIFI_BT_23 = a64.GROUND   //
	WIFI_BT_24 = a64.PL6      //
	WIFI_BT_25 = a64.VCC      //
	WIFI_BT_26 = a64.IOVCC    //

	AUDIO_LEFT  = a64.INVALID // TODO(maruel): Figure out
	AUDIO_RIGHT = a64.INVALID // TODO(maruel): Figure out
)
