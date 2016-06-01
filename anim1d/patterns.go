// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"math/rand"
	"time"
)

// Color shows a single color on all lights. It knows how to renders itself
// into a frame.
type Color struct {
	R, G, B uint8
}

// Add adds two color together with saturation, mixing according to the alpha
// channel.
func (c *Color) Add(d Color) {
	r := uint16(c.R) + uint16(d.R)
	if r > 255 {
		r = 255
	}
	c.R = uint8(r)
	g := uint16(c.G) + uint16(d.G)
	if g > 255 {
		g = 255
	}
	c.G = uint8(g)
	b := uint16(c.B) + uint16(d.B)
	if b > 255 {
		b = 255
	}
	c.B = uint8(b)
}

func (c *Color) NextFrame(pixels Frame, sinceStart time.Duration) {
	for i := range pixels {
		pixels[i] = *c
	}
}

func (c *Color) NativeDuration(pixels int) time.Duration {
	return 0
}

// Frame is a strip of colors. It knows how to renders itself into a frame
// (which is recursive).
type Frame []Color

func (f Frame) NextFrame(pixels Frame, sinceStart time.Duration) {
	copy(pixels, f)
}

func (f Frame) NativeDuration(pixels int) time.Duration {
	return 0
}

// reset() always resets the buffer to black.
func (f *Frame) reset(l int) {
	if len(*f) != l {
		*f = make(Frame, l)
	} else {
		s := *f
		for i := range s {
			s[i] = Color{}
		}
	}
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
// The trail can be a Frame or a dynamic pattern.
//
// To get smoothed movement, use Scale{} with a 5x factor or so.
// TODO(maruel): That's a bit inefficient, enable ScalingType here.
type PingPong struct {
	Trail       SPattern // [0] is the front pixel so the pixels are effectively drawn in reverse order.
	MovesPerSec float32  // Expressed in number of light jumps per second.
	buf         Frame
}

func (p *PingPong) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(pixels) == 0 || p.Trail.Pattern == nil {
		return
	}
	p.buf.reset(len(pixels)*2 - 1)
	p.Trail.NextFrame(p.buf, sinceStart)
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	//
	// For Trail being "01234567":
	//   move == 0  -> "01234567"
	//   move == 2  -> "21056789"
	//   move == 5  -> "543210ab"
	//   move == 7  -> "76543210"
	//   move == 9  -> "98765012"
	//   move == 11 -> "ba901234"
	//   move == 13 -> "d0123456"
	//   move 14 -> move 0; "2*(8-1)"
	cycle := 2 * (len(pixels) - 1)
	pos := int(float32(sinceStart.Seconds())*p.MovesPerSec) % cycle

	// It has to be copied manually because the ordering is not linear.
	// Once it works the following code looks trivial but everytime it takes me
	// an absurd amount of time to rewrite it.
	if pos >= len(pixels)-1 {
		// Runs left.
		pos = cycle - pos
		k := 2 * len(pixels)
		for i := range pixels {
			if i < pos {
				pixels[i] = p.buf[k-pos-i]
			} else {
				pixels[i] = p.buf[i-pos]
			}
		}
	} else {
		// Runs right.
		for i := range pixels {
			if i <= pos {
				pixels[i] = p.buf[pos-i]
			} else {
				pixels[i] = p.buf[pos+i]
			}
		}
	}
}

func (p *PingPong) NativeDuration(pixels int) time.Duration {
	cycle := 2 * (pixels - 1)
	return time.Duration(p.MovesPerSec*float32(cycle)) * time.Second
}

// Animation represents an animatable looping frame.
//
// If the image is smaller than the strip, doesn't touch the rest of the
// pixels. Otherwise, the excess is ignored. Use Scale{} if desired.
type Animation struct {
	Frames        []Frame
	FrameDuration time.Duration
}

func (a *Animation) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(pixels) == 0 || len(a.Frames) == 0 {
		return
	}
	copy(pixels, a.Frames[int(sinceStart/a.FrameDuration)%len(a.Frames)])
}

func (a *Animation) NativeDuration(pixels int) time.Duration {
	return a.FrameDuration * time.Duration(len(a.Frames))
}

// MakeRainbow renders rainbow colors including alpha.
type Rainbow struct {
}

func (r *Rainbow) NextFrame(pixels Frame, sinceStart time.Duration) {
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
//
// TODO(maruel): Convert to integer calculation.
func waveLength2RGB(w float32) (c Color) {
	switch {
	case 380. <= w && w < 420.:
		c.R = 128 - FloatToUint8(127.*(440.-w)/(440.-380.))
		c.B = FloatToUint8(255. * (0.1 + 0.9*(w-380.)/(420.-380.)))
	case 420. <= w && w < 440.:
		c.R = FloatToUint8(127. * (440. - w) / (440. - 380.))
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
	case 645. <= w && w < 700.:
		c.R = 255
	case 700. <= w && w < 781.:
		c.R = FloatToUint8(255. * (0.1 + 0.9*(780.-w)/(780.-700.)))
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
	Points      Frame
	MovesPerSec float32 // Expressed in number of light jumps per second.
}

func (r *Repeated) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(pixels) == 0 || len(r.Points) == 0 {
		return
	}
	offset := len(r.Points) - int(float32(sinceStart.Seconds())*r.MovesPerSec)%len(r.Points)
	for i := range pixels {
		pixels[i] = Color(r.Points[(i+offset)%len(r.Points)])
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

func (c *NightSky) NextFrame(pixels Frame, sinceStart time.Duration) {
	// random
	// animate.
}

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) NextFrame(pixels Frame, sinceStart time.Duration) {
	// TODO(maruel): Redo.
	y := float32(sinceStart.Seconds()) * 10.
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (32 + 31*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = uint8(b)
		pixels[i].B = 0
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

func (e *NightStars) NextFrame(pixels Frame, sinceStart time.Duration) {
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
				pixels[i] = Color{f, f, f}
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

func (w *WishingStar) NextFrame(pixels Frame, sinceStart time.Duration) {
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
	A, B Color
}

func (d *Gradient) NextFrame(pixels Frame, sinceStart time.Duration) {
	l := len(pixels) - 1
	if l == 0 {
		pixels[0] = d.A
		pixels[0].Add(d.B)
		return
	}
	for i := range pixels {
		/*
			// [0, 255]
			intensity := uint8(i * 255 / l)
			a := d.A
			b := d.B
			a.A = 255 - intensity
			b.A = intensity
			a.Add(b)
			pixels[i] = a
		*/
		// [0, 1]
		intensity := float32(i) / float32(len(pixels)-1)
		pixels[i] = Color{
			uint8((float32(d.A.R)*intensity + float32(d.B.R)*(1-intensity))),
			uint8((float32(d.A.G)*intensity + float32(d.B.G)*(1-intensity))),
			uint8((float32(d.A.B)*intensity + float32(d.B.B)*(1-intensity))),
		}
	}
}
