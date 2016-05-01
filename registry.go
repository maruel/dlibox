// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"sync"
	"time"
)

// PatternRegistry handles predefined patterns and their thumbnails.
type PatternRegistry struct {
	// Patterns is a map of nice predefined patterns.
	Patterns         map[string]Pattern
	NumberLEDs       int               // Must be set before calling Thumbnail().
	ThumbnailHz      int               // Must be set before calling Thumbnail().
	ThumbnailSeconds int               // Must be set before calling Thumbnail().
	cache            map[string][]byte // Thumbnail as GIF.
	lock             sync.Mutex
}

func isPixelsEqual(a, b []color.NRGBA) bool {
	if len(a) != len(b) {
		return false
	}
	for j, p := range a {
		if b[j] != p {
			return false
		}
	}
	return true
}

func (p *PatternRegistry) Thumbnail(name string) []byte {
	p.lock.Lock()
	if p.cache == nil {
		p.cache = make(map[string][]byte)
	}
	img, ok := p.cache[name]
	p.lock.Unlock()
	if ok {
		return img
	}

	pat, ok := p.Patterns[name]
	if !ok {
		return nil
	}
	pixels := [][]color.NRGBA{make([]color.NRGBA, p.NumberLEDs), make([]color.NRGBA, p.NumberLEDs)}
	nbImg := p.ThumbnailSeconds * p.ThumbnailHz
	// Change dark blue (color index #1) to background, so it can be used to save
	// more on GIF size. It's better than losing black, which is the default. To
	// not confused the Index() function, set both to the same color, so index 1
	// will never be returned by this function.
	pal := make(color.Palette, 256)
	copy(pal, palette.Plan9)
	pal[1] = pal[0]
	g := &gif.GIF{
		Image:           make([]*image.Paletted, 0, nbImg),
		Delay:           make([]int, 0, nbImg),
		Disposal:        make([]byte, 0, nbImg),
		Config:          image.Config{pal, p.NumberLEDs, 1},
		BackgroundIndex: 1,
	}
	frameDuration := int(roundF(100. / float32(p.ThumbnailHz)))
	for frame := 0; frame < nbImg; frame++ {
		since := (time.Second*time.Duration(frame) + time.Duration(p.ThumbnailHz) - 1) / time.Duration(p.ThumbnailHz)
		pat.NextFrame(pixels[frame&1], since)
		if frame > 0 && isPixelsEqual(pixels[0], pixels[1]) {
			// Skip a frame completely if its pixels didn't change at all from the
			// previous frame.
			g.Delay[len(g.Delay)-1] += frameDuration
			continue
		}
		g.Delay = append(g.Delay, frameDuration)
		g.Disposal = append(g.Disposal, gif.DisposalPrevious)
		g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, p.NumberLEDs, 1), pal))
		img := g.Image[len(g.Image)-1]
		// Compare with previous image.
		for j, p := range pixels[frame&1] {
			// TODO(maruel): draw.FloydSteinberg assuming we use 5x5 boxes instead of
			// 1x1. For now, just use the closest color.
			img.Pix[j] = uint8(pal.Index(p))
		}
	}
	b := &bytes.Buffer{}
	if err := gif.EncodeAll(b, g); err != nil {
		panic(err)
	}
	p.lock.Lock()
	p.cache[name] = b.Bytes()
	p.lock.Unlock()
	return p.cache[name]
}
