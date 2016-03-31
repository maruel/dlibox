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
}

func (p *PatternRegistry) Thumbnail(name string) []byte {
	if p.cache == nil {
		p.cache = make(map[string][]byte)
	}
	if img, ok := p.cache[name]; ok {
		return img
	}

	pat := p.Patterns[name]
	pixels := make([]color.NRGBA, p.NumberLEDs)
	nbImg := p.ThumbnailSeconds * p.ThumbnailHz
	g := &gif.GIF{Image: make([]*image.Paletted, nbImg), Delay: make([]int, nbImg)}
	d := int(roundF(100. / float64(p.ThumbnailHz)))
	for i := range g.Image {
		g.Delay[i] = d
		g.Image[i] = image.NewPaletted(image.Rect(0, 0, p.NumberLEDs, 1), palette.Plan9)
		since := (time.Second*time.Duration(i) + time.Duration(p.ThumbnailHz) - 1) / time.Duration(p.ThumbnailHz)
		pat.NextFrame(pixels, since)
		for j, p := range pixels {
			// For now, just use the closest color.
			// TODO(maruel): draw.FloydSteinberg
			g.Image[i].Pix[j] = uint8(g.Image[i].Palette.Index(p))
		}
	}
	b := &bytes.Buffer{}
	_ = gif.EncodeAll(b, g)
	p.cache[name] = b.Bytes()
	return p.cache[name]
}
