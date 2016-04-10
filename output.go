// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"bytes"
	"image/color"
	"io"
	"time"

	"github.com/maruel/ansi256"
	"github.com/mattn/go-colorable"
)

type DotStar struct {
	// Gamma correction then power limiter.
	RedGamma   float64
	RedMax     float64
	GreenGamma float64
	GreenMax   float64
	BlueGamma  float64
	BlueMax    float64
	AmpPerLED  float64
	AmpBudget  float64

	w   io.WriteCloser
	buf []byte
}

func (d *DotStar) Close() error {
	w := d.w
	d.w = nil
	return w.Close()
}

// brightnessRamp is a look up table to convert from 0-255 to 1-0x1EE1 to
// maximize the visible contrast as output by a APA102 LED.
//
// It is intentionally trying to minimize the high intensity values to lower
// Amperage.
//
// TODO(maruel): Smooth out the inflection points by defining a bezier curve,
// it's currently very rough.
var brightnessRamp = []uint32{
	0,
	1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	10,
	11,
	12,
	13,
	14,
	15,
	16,
	17,
	18,
	19,
	20,
	21,
	22,
	23,
	24,
	25,
	26,
	27,
	28,
	29,
	30,
	1 * 31, // 1/255 at full brightness.
	47,     // 1.5 * 31
	2 * 31, //
	78,     // 2.5 * 31
	3 * 31, //
	109,    // 3.5 * 31
	4 * 31,
	5 * 31,
	6 * 31,
	7 * 31,
	8 * 31,
	9 * 31,
	10 * 31,
	11 * 31,
	12 * 31,
	13 * 31,
	14 * 31,
	15 * 31,
	16 * 31,
	17 * 31,
	18 * 31,
	19 * 31,
	20 * 31,
	21 * 31,
	22 * 31,
	23 * 31,
	24 * 31,
	25 * 31,
	26 * 31,
	27 * 31,
	28 * 31,
	29 * 31,
	30 * 31,
	31 * 31,
	32 * 31,
	33 * 31,
	34 * 31,
	35 * 31,
	36 * 31,
	37 * 31,
	38 * 31,
	39 * 31,
	40 * 31,
	41 * 31,
	42 * 31,
	43 * 31,
	44 * 31,
	45 * 31,
	46 * 31,
	47 * 31,
	48 * 31,
	49 * 31,
	50 * 31,
	51 * 31,
	52 * 31,
	53 * 31,
	54 * 31,
	55 * 31,
	56 * 31,
	57 * 31,
	58 * 31,
	59 * 31,
	60 * 31,
	61 * 31,
	62 * 31,
	63 * 31,
	64 * 31,
	65 * 31,
	66 * 31,
	67 * 31,
	68 * 31,
	69 * 31,
	70 * 31,
	71 * 31,
	72 * 31,
	73 * 31,
	74 * 31,
	75 * 31,
	76 * 31,
	77 * 31,
	78 * 31,
	79 * 31,
	80 * 31,
	81 * 31,
	82 * 31,
	83 * 31,
	84 * 31,
	85 * 31,
	86 * 31,
	87 * 31,
	88 * 31,
	89 * 31,
	90 * 31,
	91 * 31,
	92 * 31,
	93 * 31,
	94 * 31,
	95 * 31,
	96 * 31,
	97 * 31,
	98 * 31,
	99 * 31,
	100 * 31,
	101 * 31,
	102 * 31,
	103 * 31,
	104 * 31,
	105 * 31,
	106 * 31,
	107 * 31,
	108 * 31,
	109 * 31,
	110 * 31,
	111 * 31,
	112 * 31,
	113 * 31,
	114 * 31,
	115 * 31,
	116 * 31,
	117 * 31,
	118 * 31,
	119 * 31,
	120 * 31,
	121 * 31,
	122 * 31,
	123 * 31,
	124 * 31,
	125 * 31,
	126 * 31,
	127 * 31,
	128 * 31,
	129 * 31,
	130 * 31,
	131 * 31,
	132 * 31,
	133 * 31,
	134 * 31,
	135 * 31,
	136 * 31,
	137 * 31,
	138 * 31,
	139 * 31,
	140 * 31,
	141 * 31,
	142 * 31,
	143 * 31,
	144 * 31,
	145 * 31,
	146 * 31,
	147 * 31,
	148 * 31,
	149 * 31,
	150 * 31,
	151 * 31,
	152 * 31,
	153 * 31,
	154 * 31,
	155 * 31,
	156 * 31,
	157 * 31,
	158 * 31,
	159 * 31,
	160 * 31,
	161 * 31,
	162 * 31,
	163 * 31,
	164 * 31,
	165 * 31,
	166 * 31,
	167 * 31,
	168 * 31,
	169 * 31,
	170 * 31,
	171 * 31,
	172 * 31,
	173 * 31,
	174 * 31,
	175 * 31,
	176 * 31,
	177 * 31,
	178 * 31,
	179 * 31,
	180 * 31,
	181 * 31,
	182 * 31,
	183 * 31,
	184 * 31,
	185 * 31,
	186 * 31,
	187 * 31,
	188 * 31,
	189 * 31,
	191 * 31,
	193 * 31,
	195 * 31,
	197 * 31,
	199 * 31,
	201 * 31,
	203 * 31,
	205 * 31,
	207 * 31,
	209 * 31,
	211 * 31,
	213 * 31,
	215 * 31,
	217 * 31,
	219 * 31,
	221 * 31,
	223 * 31,
	225 * 31,
	227 * 31,
	229 * 31,
	231 * 31,
	233 * 31,
	235 * 31,
	237 * 31,
	239 * 31,
	241 * 31,
	243 * 31,
	245 * 31,
	247 * 31,
	249 * 31,
	251 * 31,
	253 * 31,
	255 * 31, // 0x1EE1
}

// ColorToAPA102 converts a color into the 4 bytes needed to control an APA-102
// LED.
//
// The response as seen by the human eye is very non-linear. The APA-102
// provides an overall brightness PWM but it is relatively slower and results
// in human visible flicker. On the other hand the minimal color (1/255) is
// still too intense at full brightness, so for very dark color, it is worth
// using the overall brightness PWM. The goal is to use brightness!=31 as
// little as possible.
//
// Global brightness frequency is 580Hz and color frequency at 19.2kHz.
// https://cpldcpu.wordpress.com/2014/08/27/apa102/
// Both are multiplicative, so brightness@50% and color@50% means an effective
// 25% duty cycle but it is not properly distributed, which is the main problem.
//
// It is unclear to me if brightness is exactly in 1/31 increment as I don't
// have an oscilloscope to confirm. Same for color in 1/255 increment.
//
// Each channel duty cycle ramps from 100% to 1/(31*255) == 1/7905.
//
// Return brighness, blue, green, red.
func ColorToAPA102(c color.NRGBA) (byte, byte, byte, byte) {
	// Evaluate alpha.
	r, g, b, _ := c.RGBA()

	// Converts each channel from 1/0xFFFF to units of 1/0x1EE1 (1/7905) in 1/255
	// units skewed towards the low end; we want 1:1 fidelity in the low end.
	// TODO(maruel): We lose precision here, could be nice to create the ramp as
	// a function.
	r2 := brightnessRamp[r>>8]
	g2 := brightnessRamp[g>>8]
	b2 := brightnessRamp[b>>8]

	// TODO(maruel): The transition between the two modes must be seamless.
	// s1 := floatToUint8(255. * math.Pow(float64(b)/65280.*d.BlueMax, 1/d.BlueGamma))
	if r2 <= 255 && g2 <= 255 && b2 <= 255 {
		// Use lower brightness.
		return byte(0xE0 + 1), byte(b2), byte(g2), byte(r2)
	} else {
		// In this case we need to use a ramp of 255-1 even for lower colors.
		return byte(0xE0 + 31), byte(b2 / 31), byte((g2 + 0) / 31), byte((r2 + 0) / 31)
	}
}

func (d *DotStar) Write(pixels []color.NRGBA) error {
	// https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf
	numLights := len(pixels)
	// End frames are needed to be able to push enough SPI clock signals due to
	// internal half-delay of data signal from each individual LED. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	l := 4*(numLights+1) + numLights/2/8 + 1
	if len(d.buf) != l {
		d.buf = make([]byte, l)
		// Set end frames right away.
		s := d.buf[4+4*numLights:]
		for i := range s {
			s[i] = 0xFF
		}
	}
	// Start frame is all zeros. Just skip it.
	s := d.buf[4 : 4+4*numLights]
	// TODO(maruel): Calculate power in duty cycle of each channel.
	power := 0
	for i := range pixels {
		s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = ColorToAPA102(pixels[i])
		//power += p
	}
	if d.AmpBudget != 0 {
		powerF := float64(power) * d.AmpPerLED / 255.
		if powerF > d.AmpBudget {
			ratio := d.AmpBudget / powerF
			for i := range s {
				if i%4 != 0 {
					s[i] = floatToUint8(float64(s[i]) * ratio)
				}
			}
		}
	}
	_, err := d.w.Write(d.buf)
	return err
}

func (d *DotStar) MinDelay() time.Duration {
	// As per APA102-C spec, it's max refresh rate is 400hz.
	// https://en.wikipedia.org/wiki/Flicker_fusion_threshold is a recommended
	// reading.
	return time.Second / 400
}

// MakeDotStar returns a strip that communicates over SPI to APA102 LEDs.
//
// This is generally what you want once the hardware is connected.
func MakeDotStar() (*DotStar, error) {
	// The speed must be high, as there's 32 bits sent per LED, creating a
	// staggered effect. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	w, err := MakeSPI("", 10000000)
	if err != nil {
		return nil, err
	}
	return &DotStar{
		RedGamma:   1.,
		RedMax:     0.5,
		GreenGamma: 1.,
		GreenMax:   0.5,
		BlueGamma:  1.,
		BlueMax:    0.5,
		AmpPerLED:  .02,
		AmpBudget:  9.,
		w:          w,
	}, err
}

//

type screenStrip struct {
	w   io.Writer
	buf bytes.Buffer
}

func (s *screenStrip) Close() error {
	return nil
}

func (s *screenStrip) Write(pixels []color.NRGBA) error {
	// This code is designed to minimize the amount of memory allocated per call.
	s.buf.Reset()
	_, _ = s.buf.WriteString("\r\033[0m")
	for _, c := range pixels {
		_, _ = io.WriteString(&s.buf, ansi256.Default.Block(c))
	}
	_, _ = s.buf.WriteString("\033[0m ")
	_, err := s.buf.WriteTo(s.w)
	return err
}

func (s *screenStrip) MinDelay() time.Duration {
	// Limit to 30hz, especially for ssh connections.
	return time.Second / 30
}

// MakeScreen returns a strip that displays at the console.
//
// This is generally what you want while waiting for the LED strip to be
// shipped and you are excited to try it out.
func MakeScreen() Strip {
	return &screenStrip{w: colorable.NewColorableStdout()}
}

/*
type piBlasterStrip struct {
}

func (p *piBlasterStrip) Close() error {
	return nil
}

func (p *piBlasterStrip) Write(pixels []color.NRGBA) error {
	return nil
}

func (p *piBlasterStrip) MinDelay() time.Duration {
	return time.Second / 100
}

// MakePiBlaster returns a strip that control a LED through PiBlaster.
//
// This is very specific, assuming you have a 3 colors LED connected manually.
//
// This requires https://github.com/sarfata/pi-blaster to be installed.
func MakePiBlaster() Strip {
	return &piBlasterStrip{}
}
*/

//

/*
// intensityLimiter limits the maximum intensity. Does this by scaling the
// alpha channel.
type intensityLimiter struct {
	Child Pattern
	Max   int // Maximum value between 0 (off) to 255 (full intensity).
}

func (i *intensityLimiter) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	i.Child.NextFrame(pixels, sinceStart)
	for j := range pixels {
		pixels[j].A = uint8((int(pixels[j].A) + i.Max - 1) * 255 / i.Max)
	}
}

// powerLimiter limits the maximum power draw (in Amp).
//
// It does this by scaling -each- the alpha channel but only when too much LEDs
// are lit, which would cause too much Amperes to be drawn. This means when
// only a subset of the strip is lit, all colors can be used but when all the
// strip is used, the intensity is limited.
//
// TODO(maruel): Calculate the actual power draw per channel.
// TODO(maruel): Check if the draw is linear to the intensity value per channel.
// TODO(maruel): This should only be done once alpha has been evaluated.
// TODO(maruel): This shoudl only be done after gamma correction (?)
type powerLimiter struct {
	Child     Pattern
	AmpPerLED float64
	AmpBudget    float64
}

func (p *powerLimiter) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	p.Child.NextFrame(pixels, sinceStart)
	power := 0.
	for _, c := range pixels {
		cR, cG, cB, _ := c.RGBA()
		power += float64(cR>>8+cG>>8+cB>>8) * p.AmpPerLED
	}
	if power > p.AmpBudget {
		// We only need to scale down the alpha as long as we treat each channel as
		// having the same power budget.
		for i := range pixels {
			pixels[i].A = floatToUint8(float64(pixels[i].A) * power / p.AmpBudget)
		}
	}
}

// gammaCorrection corrects the intensity of each channel and 'applies' the
// alpha channel.
//
// TODO(maruel): The alpha channel should be dropped after this? As the alpha
// correction is linear.
//
// For example, the green channel will likely be much brighter than red and
// blue.
//
// '*Max' value are what should be considered 1.0, when it's deemed not
// necessary to use the channel at full intensity. This is useful as this can
// limit the amperage used by the LED strip, which is a concern for longer
// strips.
type gammaCorrection struct {
	Child      Pattern
	RedGamma   float64
	RedMax     float64
	GreenGamma float64
	GreenMax   float64
	BlueGamma  float64
	BlueMax    float64
}

func (g *gammaCorrection) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	g.Child.NextFrame(pixels, sinceStart)
	for i := range pixels {
		pixels[i].R = floatToUint8(255. * math.Pow(float64(pixels[i].R)/255.*g.RedMax, 1/g.RedGamma))
		pixels[i].G = floatToUint8(255. * math.Pow(float64(pixels[i].G)/255.*g.GreenMax, 1/g.GreenGamma))
		pixels[i].B = floatToUint8(255. * math.Pow(float64(pixels[i].B)/255.*g.BlueMax, 1/g.BlueGamma))
	}
}
*/
