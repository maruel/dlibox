// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pine64

import (
	"github.com/maruel/dlibox/go/pio/drivers"
	"github.com/maruel/dlibox/go/pio/host/allwinner"
	"github.com/maruel/dlibox/go/pio/host/allwinner_pl"
	"github.com/maruel/dlibox/go/pio/host/headers"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/protocols/analog"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

var (
	VCC         gpio.PinIO = &pins.BasicPin{Name: "VCC"}         //
	DC_IN       gpio.PinIO = &pins.BasicPin{Name: "DC_IN"}       //
	TEMP_SENSOR gpio.PinIO = &pins.BasicPin{Name: "TEMP_SENSOR"} //
	BAT_PLUS    gpio.PinIO = &pins.BasicPin{Name: "BAT_PLUS"}    //
	IR_RX       gpio.PinIO = &pins.BasicPin{Name: "IR_RX"}       // IR Data Receive
	CHARGER_LED gpio.PinIO = &pins.BasicPin{Name: "CHARGER_LED"} //
	RESET       gpio.PinIO = &pins.BasicPin{Name: "RESET"}       //
	PWR_SWITCH  gpio.PinIO = &pins.BasicPin{Name: "PWR_SWITCH"}  //
	IOVCC       gpio.PinIO = &pins.BasicPin{Name: "IOVCC"}       // Power supply for port A
)

var (
	P1_1  gpio.PinIO = pins.V3_3         // 3.3 volt; max 40mA
	P1_2  gpio.PinIO = pins.V5           // 5 volt (before filtering)
	P1_3  gpio.PinIO = allwinner.PH3     //
	P1_4  gpio.PinIO = pins.V5           //
	P1_5  gpio.PinIO = allwinner.PH2     //
	P1_6  gpio.PinIO = pins.GROUND       //
	P1_7  gpio.PinIO = allwinner_pl.PL10 //
	P1_8  gpio.PinIO = allwinner.PB0     //
	P1_9  gpio.PinIO = pins.GROUND       //
	P1_10 gpio.PinIO = allwinner.PB1     //
	P1_11 gpio.PinIO = allwinner.PC7     //
	P1_12 gpio.PinIO = allwinner.PC8     //
	P1_13 gpio.PinIO = allwinner.PH9     //
	P1_14 gpio.PinIO = pins.GROUND       //
	P1_15 gpio.PinIO = allwinner.PC12    //
	P1_16 gpio.PinIO = allwinner.PC13    //
	P1_17 gpio.PinIO = pins.V3_3         //
	P1_18 gpio.PinIO = allwinner.PC14    //
	P1_19 gpio.PinIO = allwinner.PC0     //
	P1_20 gpio.PinIO = pins.GROUND       //
	P1_21 gpio.PinIO = allwinner.PC1     //
	P1_22 gpio.PinIO = allwinner.PC15    //
	P1_23 gpio.PinIO = allwinner.PC2     //
	P1_24 gpio.PinIO = allwinner.PC3     //
	P1_25 gpio.PinIO = pins.GROUND       //
	P1_26 gpio.PinIO = allwinner.PH7     //
	P1_27 gpio.PinIO = allwinner_pl.PL9  //
	P1_28 gpio.PinIO = allwinner_pl.PL8  //
	P1_29 gpio.PinIO = allwinner.PH5     //
	P1_30 gpio.PinIO = pins.GROUND       //
	P1_31 gpio.PinIO = allwinner.PH6     //
	P1_32 gpio.PinIO = allwinner.PC4     //
	P1_33 gpio.PinIO = allwinner.PC5     //
	P1_34 gpio.PinIO = pins.GROUND       //
	P1_35 gpio.PinIO = allwinner.PC9     //
	P1_36 gpio.PinIO = allwinner.PC6     //
	P1_37 gpio.PinIO = allwinner.PC16    //
	P1_38 gpio.PinIO = allwinner.PC10    //
	P1_39 gpio.PinIO = pins.GROUND       //
	P1_40 gpio.PinIO = allwinner.PC11    //

	EULER_1  gpio.PinIO   = pins.V3_3         //
	EULER_2  gpio.PinIO   = DC_IN             //
	EULER_3  gpio.PinIO   = BAT_PLUS          //
	EULER_4  gpio.PinIO   = DC_IN             //
	EULER_5  gpio.PinIO   = TEMP_SENSOR       //
	EULER_6  gpio.PinIO   = pins.GROUND       //
	EULER_7  gpio.PinIO   = IR_RX             //
	EULER_8  gpio.PinIO   = pins.V5           //
	EULER_9  gpio.PinIO   = pins.GROUND       //
	EULER_10 gpio.PinIO   = allwinner.PH8     //
	EULER_11 gpio.PinIO   = allwinner.PB3     //
	EULER_12 gpio.PinIO   = allwinner.PB4     //
	EULER_13 gpio.PinIO   = allwinner.PB5     //
	EULER_14 gpio.PinIO   = pins.GROUND       //
	EULER_15 gpio.PinIO   = allwinner.PB6     //
	EULER_16 gpio.PinIO   = allwinner.PB7     //
	EULER_17 gpio.PinIO   = pins.V3_3         //
	EULER_18 gpio.PinIO   = allwinner.PD4     //
	EULER_19 gpio.PinIO   = allwinner.PD2     //
	EULER_20 gpio.PinIO   = pins.GROUND       //
	EULER_21 gpio.PinIO   = allwinner.PD3     //
	EULER_22 gpio.PinIO   = allwinner.PD5     //
	EULER_23 gpio.PinIO   = allwinner.PD1     //
	EULER_24 gpio.PinIO   = allwinner.PD0     //
	EULER_25 gpio.PinIO   = pins.GROUND       //
	EULER_26 gpio.PinIO   = allwinner.PD6     //
	EULER_27 gpio.PinIO   = allwinner.PB2     //
	EULER_28 gpio.PinIO   = allwinner.PD7     //
	EULER_29 gpio.PinIO   = allwinner.PB8     //
	EULER_30 gpio.PinIO   = allwinner.PB9     //
	EULER_31 analog.PinIO = allwinner.EAROUTP //
	EULER_32 analog.PinIO = allwinner.EAROUTN //
	EULER_33 gpio.PinIO   = pins.INVALID      //
	EULER_34 gpio.PinIO   = pins.GROUND       //

	EXP_1  gpio.PinIO   = pins.V3_3         //
	EXP_2  gpio.PinIO   = allwinner_pl.PL7  //
	EXP_3  gpio.PinIO   = CHARGER_LED       //
	EXP_4  gpio.PinIO   = RESET             //
	EXP_5  gpio.PinIO   = PWR_SWITCH        //
	EXP_6  gpio.PinIO   = pins.GROUND       //
	EXP_7  gpio.PinIO   = allwinner.PB8     //
	EXP_8  gpio.PinIO   = allwinner.PB9     //
	EXP_9  gpio.PinIO   = pins.GROUND       //
	EXP_10 analog.PinIO = allwinner.KEY_ADC //

	WIFI_BT_1  gpio.PinIO = pins.GROUND        //
	WIFI_BT_2  gpio.PinIO = allwinner.PG6      //
	WIFI_BT_3  gpio.PinIO = allwinner.PG0      //
	WIFI_BT_4  gpio.PinIO = allwinner.PG7      //
	WIFI_BT_5  gpio.PinIO = pins.GROUND        //
	WIFI_BT_6  gpio.PinIO = allwinner.PG8      //
	WIFI_BT_7  gpio.PinIO = allwinner.PG1      //
	WIFI_BT_8  gpio.PinIO = allwinner.PG9      //
	WIFI_BT_9  gpio.PinIO = allwinner.PG2      //
	WIFI_BT_10 gpio.PinIO = allwinner.PG10     //
	WIFI_BT_11 gpio.PinIO = allwinner.PG3      //
	WIFI_BT_12 gpio.PinIO = allwinner.PG11     //
	WIFI_BT_13 gpio.PinIO = allwinner.PG4      //
	WIFI_BT_14 gpio.PinIO = allwinner.PG12     //
	WIFI_BT_15 gpio.PinIO = allwinner.PG5      //
	WIFI_BT_16 gpio.PinIO = allwinner.PG13     //
	WIFI_BT_17 gpio.PinIO = allwinner_pl.PL2   //
	WIFI_BT_18 gpio.PinIO = pins.GROUND        //
	WIFI_BT_19 gpio.PinIO = allwinner_pl.PL3   //
	WIFI_BT_20 gpio.PinIO = allwinner_pl.PL5   //
	WIFI_BT_21 gpio.PinIO = allwinner.X32KFOUT //
	WIFI_BT_22 gpio.PinIO = allwinner_pl.PL5   //
	WIFI_BT_23 gpio.PinIO = pins.GROUND        //
	WIFI_BT_24 gpio.PinIO = allwinner_pl.PL6   //
	WIFI_BT_25 gpio.PinIO = VCC                //
	WIFI_BT_26 gpio.PinIO = IOVCC              //

	AUDIO_LEFT  gpio.PinIO = pins.INVALID // TODO(maruel): Figure out, is that EAROUT?
	AUDIO_RIGHT gpio.PinIO = pins.INVALID //
)

func zapPins() {
	P1_1 = pins.INVALID
	P1_2 = pins.INVALID
	P1_3 = pins.INVALID
	P1_4 = pins.INVALID
	P1_5 = pins.INVALID
	P1_6 = pins.INVALID
	P1_7 = pins.INVALID
	P1_8 = pins.INVALID
	P1_9 = pins.INVALID
	P1_10 = pins.INVALID
	P1_11 = pins.INVALID
	P1_12 = pins.INVALID
	P1_13 = pins.INVALID
	P1_14 = pins.INVALID
	P1_15 = pins.INVALID
	P1_16 = pins.INVALID
	P1_17 = pins.INVALID
	P1_18 = pins.INVALID
	P1_19 = pins.INVALID
	P1_20 = pins.INVALID
	P1_21 = pins.INVALID
	P1_22 = pins.INVALID
	P1_23 = pins.INVALID
	P1_24 = pins.INVALID
	P1_25 = pins.INVALID
	P1_26 = pins.INVALID
	P1_27 = pins.INVALID
	P1_28 = pins.INVALID
	P1_29 = pins.INVALID
	P1_30 = pins.INVALID
	P1_31 = pins.INVALID
	P1_32 = pins.INVALID
	P1_33 = pins.INVALID
	P1_34 = pins.INVALID
	P1_35 = pins.INVALID
	P1_36 = pins.INVALID
	P1_37 = pins.INVALID
	P1_38 = pins.INVALID
	P1_39 = pins.INVALID
	P1_40 = pins.INVALID
	EULER_1 = pins.INVALID
	EULER_2 = pins.INVALID
	EULER_3 = pins.INVALID
	EULER_4 = pins.INVALID
	EULER_5 = pins.INVALID
	EULER_6 = pins.INVALID
	EULER_7 = pins.INVALID
	EULER_8 = pins.INVALID
	EULER_9 = pins.INVALID
	EULER_10 = pins.INVALID
	EULER_11 = pins.INVALID
	EULER_12 = pins.INVALID
	EULER_13 = pins.INVALID
	EULER_14 = pins.INVALID
	EULER_15 = pins.INVALID
	EULER_16 = pins.INVALID
	EULER_17 = pins.INVALID
	EULER_18 = pins.INVALID
	EULER_19 = pins.INVALID
	EULER_20 = pins.INVALID
	EULER_21 = pins.INVALID
	EULER_22 = pins.INVALID
	EULER_23 = pins.INVALID
	EULER_24 = pins.INVALID
	EULER_25 = pins.INVALID
	EULER_26 = pins.INVALID
	EULER_27 = pins.INVALID
	EULER_28 = pins.INVALID
	EULER_29 = pins.INVALID
	EULER_30 = pins.INVALID
	EULER_31 = pins.INVALID
	EULER_32 = pins.INVALID
	EULER_33 = pins.INVALID
	EULER_34 = pins.INVALID
	EXP_1 = pins.INVALID
	EXP_2 = pins.INVALID
	EXP_3 = pins.INVALID
	EXP_4 = pins.INVALID
	EXP_5 = pins.INVALID
	EXP_6 = pins.INVALID
	EXP_7 = pins.INVALID
	EXP_8 = pins.INVALID
	EXP_9 = pins.INVALID
	EXP_10 = pins.INVALID
	WIFI_BT_1 = pins.INVALID
	WIFI_BT_2 = pins.INVALID
	WIFI_BT_3 = pins.INVALID
	WIFI_BT_4 = pins.INVALID
	WIFI_BT_5 = pins.INVALID
	WIFI_BT_6 = pins.INVALID
	WIFI_BT_7 = pins.INVALID
	WIFI_BT_8 = pins.INVALID
	WIFI_BT_9 = pins.INVALID
	WIFI_BT_10 = pins.INVALID
	WIFI_BT_11 = pins.INVALID
	WIFI_BT_12 = pins.INVALID
	WIFI_BT_13 = pins.INVALID
	WIFI_BT_14 = pins.INVALID
	WIFI_BT_15 = pins.INVALID
	WIFI_BT_16 = pins.INVALID
	WIFI_BT_17 = pins.INVALID
	WIFI_BT_18 = pins.INVALID
	WIFI_BT_19 = pins.INVALID
	WIFI_BT_20 = pins.INVALID
	WIFI_BT_21 = pins.INVALID
	WIFI_BT_22 = pins.INVALID
	WIFI_BT_23 = pins.INVALID
	WIFI_BT_24 = pins.INVALID
	WIFI_BT_25 = pins.INVALID
	WIFI_BT_26 = pins.INVALID
}

// driver implements drivers.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "pine64"
}

func (d *driver) Type() drivers.Type {
	return drivers.Pins
}

func (d *driver) Prerequisites() []string {
	return []string{"allwinner_pl"}
}

func (d *driver) Init() (bool, error) {
	if !internal.IsPine64() {
		zapPins()
		return false, nil
	}
	headers.Register("P1", [][]pins.Pin{
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
	})
	headers.Register("EULER", [][]pins.Pin{
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
	})
	headers.Register("EXP", [][]pins.Pin{
		{EXP_1, EXP_2},
		{EXP_3, EXP_4},
		{EXP_5, EXP_6},
		{EXP_7, EXP_8},
		{EXP_9, EXP_10},
	})
	headers.Register("WIFI_BT", [][]pins.Pin{
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
	})
	headers.Register("AUDIO", [][]pins.Pin{
		{AUDIO_LEFT},
		{AUDIO_RIGHT},
	})
	return true, nil
}

var _ drivers.Driver = &driver{}
