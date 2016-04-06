// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"bytes"
	"image/color"
	"io"
	"math"
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
	// Use a brightness constant of 31 since the brightness PWM is poor on APA102.
	brightness := byte(0xE0 + 31)
	power := 0
	for i := range pixels {
		// Evaluate alpha.
		r, g, b, _ := pixels[i].RGBA()
		// BGR.
		s1 := floatToUint8(255. * math.Pow(float64(b)/65280.*d.BlueMax, 1/d.BlueGamma))
		s2 := floatToUint8(255. * math.Pow(float64(g)/65280.*d.GreenMax, 1/d.GreenGamma))
		s3 := floatToUint8(255. * math.Pow(float64(r)/65280.*d.RedMax, 1/d.RedGamma))
		power += int(s1) + int(s2) + int(s3)
		s[4*i] = brightness
		s[4*i+1] = s1
		s[4*i+2] = s2
		s[4*i+3] = s3
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
	w, err := makeSPI("", 10000000)
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
