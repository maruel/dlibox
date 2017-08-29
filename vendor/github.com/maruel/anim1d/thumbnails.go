// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"runtime"
	"sync"
)

// ThumbnailsCache is a cache of animated GIF thumbnails for each pattern.
type ThumbnailsCache struct {
	NumberLEDs       int // Must be set before calling Thumbnail().
	ThumbnailHz      int // Must be set before calling Thumbnail().
	ThumbnailSeconds int // Must be set before calling Thumbnail().

	lock  sync.Mutex
	c     chan struct{}     // Limits the number of concurrent GIF animation to number of CPU core.
	cache map[string][]byte // Thumbnail as GIF. The key is the JSON serialized form encoded as a string.
}

// GIF returns a serialized animated GIF for a JSON serialized pattern.
func (t *ThumbnailsCache) GIF(serialized []byte) ([]byte, error) {
	k := string(serialized)

	t.lock.Lock()
	if t.cache == nil {
		// Fresh object, initialize it.
		t.c = make(chan struct{}, runtime.NumCPU())
		for n := 0; n < cap(t.c); n++ {
			t.c <- struct{}{}
		}
		t.cache = make(map[string][]byte)
	}
	img, ok := t.cache[k]
	t.lock.Unlock()

	if ok {
		return img, nil
	}

	// Limit number of operations to number of CPU cores. Particularly important
	// on slower platform.
	<-t.c
	defer func() {
		t.c <- struct{}{}
	}()

	// Unmarshal the string to recreate the Pattern object.
	var pat SPattern
	if err := json.Unmarshal(serialized, &pat); err != nil {
		return nil, err
	}
	pixels := []Frame{make(Frame, t.NumberLEDs), make(Frame, t.NumberLEDs)}
	nbImg := t.ThumbnailSeconds * t.ThumbnailHz
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
		Config:          image.Config{ColorModel: pal, Width: t.NumberLEDs, Height: 1},
		BackgroundIndex: 1,
	}
	frameDuration := (100 + t.ThumbnailHz>>1) / t.ThumbnailHz
	for frame := 0; frame < nbImg; frame++ {
		since := uint32(1000 * frame / t.ThumbnailHz)
		pat.Render(pixels[frame&1], since)
		if frame > 0 && pixels[0].isEqual(pixels[1]) {
			// Skip a frame completely if its pixels didn't change at all from the
			// previous frame.
			g.Delay[len(g.Delay)-1] += frameDuration
			continue
		}
		g.Delay = append(g.Delay, frameDuration)
		g.Disposal = append(g.Disposal, gif.DisposalPrevious)
		g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, t.NumberLEDs, 1), pal))
		img := g.Image[len(g.Image)-1]
		// Compare with previous image.
		for j, pixel := range pixels[frame&1] {
			// TODO(maruel): draw.FloydSteinberg assuming we use 5x5 boxes instead of
			// 1x1. For now, just use the closest color.
			c := color.NRGBA{pixel.R, pixel.G, pixel.B, 255}
			img.Pix[j] = uint8(pal.Index(c))
		}
	}
	b := &bytes.Buffer{}
	if err := gif.EncodeAll(b, g); err != nil {
		panic(err)
	}
	out := b.Bytes()

	t.lock.Lock()
	t.cache[k] = out
	t.lock.Unlock()

	return out, nil
}
