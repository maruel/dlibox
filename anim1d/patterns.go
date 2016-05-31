// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"image/color"
	"image/png"
	"math/rand"
	"time"
)

// RainbowColors are approximate rainbow colors without alpha.
var RainbowColors = []color.NRGBA{
	{255, 0, 0, 255},
	{255, 127, 0, 255},
	{255, 255, 0, 255},
	{0, 255, 0, 255},
	{0, 0, 255, 255},
	{75, 0, 130, 255},
	{139, 0, 255, 255},
}

// K2000Colors can be used with PingPong to look like Knight Rider.
// https://en.wikipedia.org/wiki/Knight_Rider_(1982_TV_series)
var K2000Colors = []color.NRGBA{
	{0xff, 0, 0, 255},
	{0xff, 0, 0, 255},
	{0xee, 0, 0, 255},
	{0xdd, 0, 0, 255},
	{0xcc, 0, 0, 255},
	{0xbb, 0, 0, 255},
	{0xaa, 0, 0, 255},
	{0x99, 0, 0, 255},
	{0x88, 0, 0, 255},
	{0x77, 0, 0, 255},
	{0x66, 0, 0, 255},
	{0x55, 0, 0, 255},
	{0x44, 0, 0, 255},
	{0x33, 0, 0, 255},
	{0x22, 0, 0, 255},
	{0x11, 0, 0, 255},
}

// Color shows a single color on all lights.
type Color color.NRGBA

func (c *Color) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	cc := color.NRGBA(*c)
	for i := range pixels {
		pixels[i] = cc
	}
}

func (c *Color) NativeDuration(pixels int) time.Duration {
	return 0
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
type PingPong struct {
	Trail       []color.NRGBA // [0] is the front pixel.
	Background  color.NRGBA
	MovesPerSec float32 // Expressed in number of light jumps per second.
}

func (p *PingPong) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	for i := range pixels {
		pixels[i] = p.Background
	}
	if len(pixels) < 2 || len(p.Trail) == 0 {
		// Not worth special casing for len(pixels)==1.
		return
	}
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	cycle := 2 * (len(pixels) - 1)
	moves := int(float32(sinceStart.Seconds()) * p.MovesPerSec)
	if moves < len(pixels) {
		// On the first cycle, the trail has not bounced yet.
		pos := moves % cycle
		for i := 0; i < len(p.Trail) && pos-i >= 0; i++ {
			pixels[pos-i] = p.Trail[i]
		}
	} else {
		for i := len(p.Trail) - 1; i >= 0; i-- {
			r := (moves - i) % cycle
			if r >= len(pixels) {
				r = cycle - r
			}
			pixels[r] = p.Trail[i]
		}
	}
}

func (p *PingPong) NativeDuration(pixels int) time.Duration {
	// TODO(maruel): Rounding.
	pixels += 2
	return time.Duration(p.MovesPerSec*float32(pixels)) * time.Second
}

// Animation represents an animatable looping frame.
//
// If the image is smaller than the strip, doesn't touch the rest of the
// pixels. Otherwise, the excess is ignored. Use Scale{} if desired.
type Animation struct {
	Frames        [][]color.NRGBA
	FrameDuration time.Duration
}

// LoadAnimate loads an Animation from a PNG file.
//
// Returns nil if the file can't be found. If vertical is true, rotate the
// image by 90°.
func LoadAnimate(content []byte, frameDuration time.Duration, vertical bool) *Animation {
	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	bounds := img.Bounds()
	maxY := bounds.Max.Y
	maxX := bounds.Max.X
	if vertical {
		// Invert axes.
		maxY, maxX = maxX, maxY
	}
	buf := make([][]color.NRGBA, maxY)
	for y := 0; y < maxY; y++ {
		buf[y] = make([]color.NRGBA, maxX)
		for x := 0; x < maxX; x++ {
			if vertical {
				buf[y][x] = color.NRGBAModel.Convert(img.At(y, x)).(color.NRGBA)
			} else {
				buf[y][x] = color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			}
		}
	}
	return &Animation{buf, frameDuration}
}

func (a *Animation) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(pixels) == 0 || len(a.Frames) == 0 {
		return
	}
	copy(pixels, a.Frames[int(sinceStart/a.FrameDuration)%len(a.Frames)])
}

func (a *Animation) NativeDuration(pixels int) time.Duration {
	return a.FrameDuration * time.Duration(len(a.Frames))
}

// MakeRainbow returns rainbow colors including alpha.
type Rainbow struct {
}

func (r *Rainbow) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	start := float32(380.)
	end := float32(781.)
	/*
		step := (end - start) / float32(len(pixels)-1)
			for i := range pixels {
				// TODO(maruel): Use log scale.
				pixels[i] = waveLength2RGB(start + step*float32(i))
			}
	*/

	// TODO(maruel): Still too much red not enough pink.
	delta := end - start
	scale := logn(2)
	step := 1. / float32(len(pixels))
	for i := range pixels {
		j := log1p(float32(len(pixels)-i-1)*step) / scale
		pixels[i] = waveLength2RGB(start + delta*(1-j))
	}
}

func (r *Rainbow) NativeDuration(pixels int) time.Duration {
	return 0
}

// waveLengthToRGB returns a color over a rainbow, including alpha.
//
// This code was inspired by public domain code on the internet.
func waveLength2RGB(w float32) (c color.NRGBA) {
	switch {
	case 380. <= w && w < 440.:
		c.R = FloatToUint8(255. * (440. - w) / (440. - 380.))
		c.B = 255
	case 440. <= w && w < 490.:
		c.G = FloatToUint8(255. * (w - 440.) / (490. - 440.))
		c.B = 255
	case 490. <= w && w < 510.:
		c.G = 255
		c.B = FloatToUint8(255. * (510. - w) / (510. - 490.))
	case 510. <= w && w < 580.:
		c.R = FloatToUint8(255. * (w - 510.) / (580. - 510.))
		c.G = 255
	case 580. <= w && w < 645.:
		c.R = 255
		c.G = FloatToUint8(255. * (645. - w) / (645. - 580.))
	case 645. <= w && w < 781.:
		c.R = 255
	}
	switch {
	case 380. <= w && w < 420.:
		c.A = FloatToUint8(255. * (0.1 + 0.9*(w-380.)/(420.-380.)))
	case 420. <= w && w < 701.:
		c.A = 255
	case 701. <= w && w < 781.:
		c.A = FloatToUint8(255. * (0.1 + 0.9*(780.-w)/(780.-700.)))
	}
	return
}

// Repeated prints a repeated pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
//
// Using one point results in the same as Color{}.
//
// TODO(maruel): Refactor MovesPerSec to a new mixer 'Markee'.
type Repeated struct {
	Points      []color.NRGBA
	MovesPerSec float32 // Expressed in number of light jumps per second.
}

func (r *Repeated) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(pixels) == 0 || len(r.Points) == 0 {
		return
	}
	offset := len(r.Points) - int(float32(sinceStart.Seconds())*r.MovesPerSec)%len(r.Points)
	for i := range pixels {
		pixels[i] = r.Points[(i+offset)%len(r.Points)]
	}
}

func (r *Repeated) NativeDuration(pixels int) time.Duration {
	// TODO(maruel): Rounding.
	return time.Duration(float32(len(r.Points))/r.MovesPerSec) * time.Second
}

type point struct {
	star  int
	start time.Time
}

// NightSky has:
//    - Stars
//    - WishingStar
//    - Aurores
//    - Super nova.
//    - Rotation de la terre?
//    - Station Internationale?
type NightSky struct {
	Stars     []Animation
	Frequency float32 // Number of explosions by second.
	points    []point
}

func (c *NightSky) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	// random
	// animate.
}

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	// TODO(maruel): Redo.
	y := float32(sinceStart.Seconds()) * 10.
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (32 + 31*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = 255
		pixels[i].B = 0
		pixels[i].A = uint8(b)
	}
}

type NightStar struct {
	Intensity uint8
	Type      int
}

type NightStars struct {
	Stars []NightStar
	Seed  int // Change it to create a different pseudo-random animation.
	r     *rand.Rand
}

func (e *NightStars) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if e.r == nil {
		e.r = rand.New(rand.NewSource(int64(e.Seed)))
	}
	if len(e.Stars) != len(pixels) {
		e.Stars = make([]NightStar, len(pixels))
		for i := 0; i < len(pixels); {
			// Add a star. Decide it's relative position, intensity and type.
			// ExpFloat64() ?
			f := abs(3 * float32(e.r.NormFloat64()))
			if f < 1 {
				continue
			}
			i += int(roundF(f))
			if i >= len(pixels) {
				break
			}
			// e.r.Intn(255)
			intensity := abs(float32(e.r.NormFloat64()))
			if intensity > 255 {
				intensity = 0
			}
			e.Stars[i].Intensity = FloatToUint8(intensity)
		}
	}
	for i := range e.Stars {
		if j := e.Stars[i].Intensity; j != 0 {
			// TODO(maruel): Type, oscillation.
			if j != 0 {
				f := FloatToUint8(float32(e.r.NormFloat64())*4 + float32(j))
				pixels[i] = color.NRGBA{255, 255, 255, f}
			}
		}
	}
}

// WishingStar draws a wishing star from time to time.
//
// It will only draw one star at a time. To increase the likelihood of getting
// many simultaneously, create multiple instances and use Mixer with Weights of
// 1.
type WishingStar struct {
	Duration     time.Duration // Average duration of a star.
	AverageDelay time.Duration // Average delay between each wishing star.
}

func (w *WishingStar) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	/*
		// Create a deterministic replay by using the current number of
		// the wishing star as the seed for the current flow. Make it independent of
		// any other non-deterministic source.
		i := sinceStart / w.AverageDelay
		r := rand.New(rand.NewSource(int64(i)))
		// Always calculate things in the same order to keep the calculation
		// deterministic.
		startOffset := r.Int63()
		startPos := r.Int63()
		intensity := r.Int63()
		orientation := r.Intn(2)
		// Draw according to these parameters.
		// - Trail
		// - Observed speed based on orientation
	*/
}

// Gradient does a gradient between 2 colors as a static image.
//
// TODO(maruel): Support N colors at M positions. Only support linear gradient?
type Gradient struct {
	A, B color.NRGBA
}

func (d *Gradient) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	for i := range pixels {
		// [0, 1]
		intensity := float32(i) / float32(len(pixels)-1)
		pixels[i] = color.NRGBA{
			uint8((float32(d.A.R)*intensity + float32(d.B.R)*(1-intensity))),
			uint8((float32(d.A.G)*intensity + float32(d.B.G)*(1-intensity))),
			uint8((float32(d.A.B)*intensity + float32(d.B.B)*(1-intensity))),
			uint8((float32(d.A.A)*intensity + float32(d.B.A)*(1-intensity))),
		}
	}
}
