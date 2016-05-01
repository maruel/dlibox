// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"image/color"
	"testing"
	"time"

	"github.com/maruel/ut"
)

func TestStaticColor(t *testing.T) {
	p := &StaticColor{color.NRGBA{255, 255, 255, 255}}
	e := []expectation{{3 * time.Second, []color.NRGBA{{255, 255, 255, 255}}}}
	frames(t, p, e)
}

func TestGlow1(t *testing.T) {
	p := &Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}, 1}
	e := []expectation{
		{0, []color.NRGBA{{0xFF, 0xFF, 0xFF, 0xFF}}},
		{250 * time.Millisecond, []color.NRGBA{{0xBF, 0xBF, 0xBF, 0xFF}}},
		{500 * time.Millisecond, []color.NRGBA{{0x7F, 0x7F, 0x7F, 0xFF}}},
		{750 * time.Millisecond, []color.NRGBA{{0x3F, 0x3F, 0x3F, 0xFF}}},
		{1000 * time.Millisecond, []color.NRGBA{{0x00, 0x00, 0x00, 0xFF}}},
	}
	frames(t, p, e)
}

func TestGlow2(t *testing.T) {
	p := &Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}, 0.1}
	e := []expectation{
		{0, []color.NRGBA{{0xFF, 0xFF, 0xFF, 0xFF}}},
		{2500 * time.Millisecond, []color.NRGBA{{0xBF, 0xBF, 0xBF, 0xFF}}},
		{5000 * time.Millisecond, []color.NRGBA{{0x80, 0x80, 0x80, 0xFF}}},
		{7500 * time.Millisecond, []color.NRGBA{{0x3F, 0x3F, 0x3F, 0xFF}}},
		{10000 * time.Millisecond, []color.NRGBA{{0x00, 0x00, 0x00, 0xFF}}},
	}
	frames(t, p, e)
}

func TestPingPong(t *testing.T) {
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	p := &PingPong{[]color.NRGBA{a, b}, color.NRGBA{}, 1000}
	e := []expectation{
		{0, []color.NRGBA{a, {}, {}}},
		{500 * time.Microsecond, []color.NRGBA{a, {}, {}}},
		{1 * time.Millisecond, []color.NRGBA{b, a, {}}},
		{2 * time.Millisecond, []color.NRGBA{{}, b, a}},
		{3 * time.Millisecond, []color.NRGBA{{}, a, b}},
		{4 * time.Millisecond, []color.NRGBA{a, b, {}}},
		{5 * time.Millisecond, []color.NRGBA{b, a, {}}},
		{6 * time.Millisecond, []color.NRGBA{{}, b, a}},
	}
	frames(t, p, e)
}

func TestRepeated(t *testing.T) {
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	c := color.NRGBA{30, 30, 30, 30}
	p := &Repeated{[]color.NRGBA{a, b, c}, 1000}
	e := []expectation{
		{0, []color.NRGBA{a, b, c, a, b}},
		{500 * time.Microsecond, []color.NRGBA{a, b, c, a, b}},
		{1 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{2 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{3 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
		{4 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{5 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{6 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
	}
	frames(t, p, e)
}

func TestRepeatedRev(t *testing.T) {
	// Works in reverse too.
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	c := color.NRGBA{30, 30, 30, 30}
	p := &Repeated{[]color.NRGBA{a, b, c}, -1000}
	e := []expectation{
		{0, []color.NRGBA{a, b, c, a, b}},
		{500 * time.Microsecond, []color.NRGBA{a, b, c, a, b}},
		{1 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{2 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{3 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
		{4 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{5 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{6 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
	}
	frames(t, p, e)
}

func TestWaveLength2RGB(t *testing.T) {
	data := []struct {
		input    float32
		expected color.NRGBA
	}{
		{379, color.NRGBA{0x00, 0x00, 0x00, 0x00}},
		{380, color.NRGBA{0xFF, 0x00, 0xFF, 0x1A}},
		{381, color.NRGBA{0xFC, 0x00, 0xFF, 0x20}},
		{382, color.NRGBA{0xF7, 0x00, 0xFF, 0x26}},
		{383, color.NRGBA{0xF3, 0x00, 0xFF, 0x2C}},
		{384, color.NRGBA{0xEF, 0x00, 0xFF, 0x31}},
		{385, color.NRGBA{0xEB, 0x00, 0xFF, 0x37}},
		{386, color.NRGBA{0xE6, 0x00, 0xFF, 0x3D}},
		{387, color.NRGBA{0xE2, 0x00, 0xFF, 0x43}},
		{388, color.NRGBA{0xDE, 0x00, 0xFF, 0x48}},
		{389, color.NRGBA{0xDA, 0x00, 0xFF, 0x4E}},
		{390, color.NRGBA{0xD5, 0x00, 0xFF, 0x54}},
		{391, color.NRGBA{0xD1, 0x00, 0xFF, 0x5A}},
		{392, color.NRGBA{0xCD, 0x00, 0xFF, 0x5F}},
		{393, color.NRGBA{0xC9, 0x00, 0xFF, 0x65}},
		{394, color.NRGBA{0xC4, 0x00, 0xFF, 0x6B}},
		{395, color.NRGBA{0xC0, 0x00, 0xFF, 0x71}},
		{396, color.NRGBA{0xBC, 0x00, 0xFF, 0x76}},
		{397, color.NRGBA{0xB8, 0x00, 0xFF, 0x7C}},
		{398, color.NRGBA{0xB3, 0x00, 0xFF, 0x82}},
		{399, color.NRGBA{0xAF, 0x00, 0xFF, 0x88}},
		{400, color.NRGBA{0xAB, 0x00, 0xFF, 0x8D}},
		{401, color.NRGBA{0xA7, 0x00, 0xFF, 0x93}},
		{402, color.NRGBA{0xA2, 0x00, 0xFF, 0x99}},
		{403, color.NRGBA{0x9E, 0x00, 0xFF, 0x9E}},
		{404, color.NRGBA{0x9A, 0x00, 0xFF, 0xA4}},
		{405, color.NRGBA{0x96, 0x00, 0xFF, 0xAA}},
		{406, color.NRGBA{0x91, 0x00, 0xFF, 0xB0}},
		{407, color.NRGBA{0x8D, 0x00, 0xFF, 0xB5}},
		{408, color.NRGBA{0x89, 0x00, 0xFF, 0xBB}},
		{409, color.NRGBA{0x85, 0x00, 0xFF, 0xC1}},
		{410, color.NRGBA{0x80, 0x00, 0xFF, 0xC7}},
		{411, color.NRGBA{0x7C, 0x00, 0xFF, 0xCC}},
		{412, color.NRGBA{0x78, 0x00, 0xFF, 0xD2}},
		{413, color.NRGBA{0x74, 0x00, 0xFF, 0xD8}},
		{414, color.NRGBA{0x6F, 0x00, 0xFF, 0xDE}},
		{415, color.NRGBA{0x6B, 0x00, 0xFF, 0xE3}},
		{416, color.NRGBA{0x67, 0x00, 0xFF, 0xE9}},
		{417, color.NRGBA{0x63, 0x00, 0xFF, 0xEF}},
		{418, color.NRGBA{0x5E, 0x00, 0xFF, 0xF5}},
		{419, color.NRGBA{0x5A, 0x00, 0xFF, 0xFA}},
		{420, color.NRGBA{0x56, 0x00, 0xFF, 0xFF}},
		{421, color.NRGBA{0x52, 0x00, 0xFF, 0xFF}},
		{422, color.NRGBA{0x4D, 0x00, 0xFF, 0xFF}},
		{423, color.NRGBA{0x49, 0x00, 0xFF, 0xFF}},
		{424, color.NRGBA{0x45, 0x00, 0xFF, 0xFF}},
		{425, color.NRGBA{0x41, 0x00, 0xFF, 0xFF}},
		{426, color.NRGBA{0x3C, 0x00, 0xFF, 0xFF}},
		{427, color.NRGBA{0x38, 0x00, 0xFF, 0xFF}},
		{428, color.NRGBA{0x34, 0x00, 0xFF, 0xFF}},
		{429, color.NRGBA{0x30, 0x00, 0xFF, 0xFF}},
		{430, color.NRGBA{0x2B, 0x00, 0xFF, 0xFF}},
		{431, color.NRGBA{0x27, 0x00, 0xFF, 0xFF}},
		{432, color.NRGBA{0x23, 0x00, 0xFF, 0xFF}},
		{433, color.NRGBA{0x1F, 0x00, 0xFF, 0xFF}},
		{434, color.NRGBA{0x1A, 0x00, 0xFF, 0xFF}},
		{435, color.NRGBA{0x16, 0x00, 0xFF, 0xFF}},
		{436, color.NRGBA{0x12, 0x00, 0xFF, 0xFF}},
		{437, color.NRGBA{0x0E, 0x00, 0xFF, 0xFF}},
		{438, color.NRGBA{0x09, 0x00, 0xFF, 0xFF}},
		{439, color.NRGBA{0x05, 0x00, 0xFF, 0xFF}},
		{440, color.NRGBA{0x00, 0x00, 0xFF, 0xFF}},
		{441, color.NRGBA{0x00, 0x06, 0xFF, 0xFF}},
		{442, color.NRGBA{0x00, 0x0B, 0xFF, 0xFF}},
		{443, color.NRGBA{0x00, 0x10, 0xFF, 0xFF}},
		{444, color.NRGBA{0x00, 0x15, 0xFF, 0xFF}},
		{445, color.NRGBA{0x00, 0x1A, 0xFF, 0xFF}},
		{446, color.NRGBA{0x00, 0x20, 0xFF, 0xFF}},
		{447, color.NRGBA{0x00, 0x25, 0xFF, 0xFF}},
		{448, color.NRGBA{0x00, 0x2A, 0xFF, 0xFF}},
		{449, color.NRGBA{0x00, 0x2F, 0xFF, 0xFF}},
		{450, color.NRGBA{0x00, 0x34, 0xFF, 0xFF}},
		{451, color.NRGBA{0x00, 0x39, 0xFF, 0xFF}},
		{452, color.NRGBA{0x00, 0x3E, 0xFF, 0xFF}},
		{453, color.NRGBA{0x00, 0x43, 0xFF, 0xFF}},
		{454, color.NRGBA{0x00, 0x48, 0xFF, 0xFF}},
		{455, color.NRGBA{0x00, 0x4D, 0xFF, 0xFF}},
		{456, color.NRGBA{0x00, 0x53, 0xFF, 0xFF}},
		{457, color.NRGBA{0x00, 0x58, 0xFF, 0xFF}},
		{458, color.NRGBA{0x00, 0x5D, 0xFF, 0xFF}},
		{459, color.NRGBA{0x00, 0x62, 0xFF, 0xFF}},
		{460, color.NRGBA{0x00, 0x67, 0xFF, 0xFF}},
		{461, color.NRGBA{0x00, 0x6C, 0xFF, 0xFF}},
		{462, color.NRGBA{0x00, 0x71, 0xFF, 0xFF}},
		{463, color.NRGBA{0x00, 0x76, 0xFF, 0xFF}},
		{464, color.NRGBA{0x00, 0x7B, 0xFF, 0xFF}},
		{465, color.NRGBA{0x00, 0x80, 0xFF, 0xFF}},
		{466, color.NRGBA{0x00, 0x86, 0xFF, 0xFF}},
		{467, color.NRGBA{0x00, 0x8B, 0xFF, 0xFF}},
		{468, color.NRGBA{0x00, 0x90, 0xFF, 0xFF}},
		{469, color.NRGBA{0x00, 0x95, 0xFF, 0xFF}},
		{470, color.NRGBA{0x00, 0x9A, 0xFF, 0xFF}},
		{471, color.NRGBA{0x00, 0x9F, 0xFF, 0xFF}},
		{472, color.NRGBA{0x00, 0xA4, 0xFF, 0xFF}},
		{473, color.NRGBA{0x00, 0xA9, 0xFF, 0xFF}},
		{474, color.NRGBA{0x00, 0xAE, 0xFF, 0xFF}},
		{475, color.NRGBA{0x00, 0xB3, 0xFF, 0xFF}},
		{476, color.NRGBA{0x00, 0xB9, 0xFF, 0xFF}},
		{477, color.NRGBA{0x00, 0xBE, 0xFF, 0xFF}},
		{478, color.NRGBA{0x00, 0xC3, 0xFF, 0xFF}},
		{479, color.NRGBA{0x00, 0xC8, 0xFF, 0xFF}},
		{480, color.NRGBA{0x00, 0xCD, 0xFF, 0xFF}},
		{481, color.NRGBA{0x00, 0xD2, 0xFF, 0xFF}},
		{482, color.NRGBA{0x00, 0xD7, 0xFF, 0xFF}},
		{483, color.NRGBA{0x00, 0xDC, 0xFF, 0xFF}},
		{484, color.NRGBA{0x00, 0xE1, 0xFF, 0xFF}},
		{485, color.NRGBA{0x00, 0xE6, 0xFF, 0xFF}},
		{486, color.NRGBA{0x00, 0xEC, 0xFF, 0xFF}},
		{487, color.NRGBA{0x00, 0xF1, 0xFF, 0xFF}},
		{488, color.NRGBA{0x00, 0xF6, 0xFF, 0xFF}},
		{489, color.NRGBA{0x00, 0xFB, 0xFF, 0xFF}},
		{490, color.NRGBA{0x00, 0xFF, 0xFF, 0xFF}},
		{491, color.NRGBA{0x00, 0xFF, 0xF3, 0xFF}},
		{492, color.NRGBA{0x00, 0xFF, 0xE6, 0xFF}},
		{493, color.NRGBA{0x00, 0xFF, 0xDA, 0xFF}},
		{494, color.NRGBA{0x00, 0xFF, 0xCD, 0xFF}},
		{495, color.NRGBA{0x00, 0xFF, 0xC0, 0xFF}},
		{496, color.NRGBA{0x00, 0xFF, 0xB3, 0xFF}},
		{497, color.NRGBA{0x00, 0xFF, 0xA7, 0xFF}},
		{498, color.NRGBA{0x00, 0xFF, 0x9A, 0xFF}},
		{499, color.NRGBA{0x00, 0xFF, 0x8D, 0xFF}},
		{500, color.NRGBA{0x00, 0xFF, 0x80, 0xFF}},
		{501, color.NRGBA{0x00, 0xFF, 0x74, 0xFF}},
		{502, color.NRGBA{0x00, 0xFF, 0x67, 0xFF}},
		{503, color.NRGBA{0x00, 0xFF, 0x5A, 0xFF}},
		{504, color.NRGBA{0x00, 0xFF, 0x4D, 0xFF}},
		{505, color.NRGBA{0x00, 0xFF, 0x41, 0xFF}},
		{506, color.NRGBA{0x00, 0xFF, 0x34, 0xFF}},
		{507, color.NRGBA{0x00, 0xFF, 0x27, 0xFF}},
		{508, color.NRGBA{0x00, 0xFF, 0x1A, 0xFF}},
		{509, color.NRGBA{0x00, 0xFF, 0x0E, 0xFF}},
		{510, color.NRGBA{0x00, 0xFF, 0x00, 0xFF}},
		{511, color.NRGBA{0x05, 0xFF, 0x00, 0xFF}},
		{512, color.NRGBA{0x08, 0xFF, 0x00, 0xFF}},
		{513, color.NRGBA{0x0C, 0xFF, 0x00, 0xFF}},
		{514, color.NRGBA{0x10, 0xFF, 0x00, 0xFF}},
		{515, color.NRGBA{0x13, 0xFF, 0x00, 0xFF}},
		{516, color.NRGBA{0x17, 0xFF, 0x00, 0xFF}},
		{517, color.NRGBA{0x1A, 0xFF, 0x00, 0xFF}},
		{518, color.NRGBA{0x1E, 0xFF, 0x00, 0xFF}},
		{519, color.NRGBA{0x22, 0xFF, 0x00, 0xFF}},
		{520, color.NRGBA{0x25, 0xFF, 0x00, 0xFF}},
		{521, color.NRGBA{0x29, 0xFF, 0x00, 0xFF}},
		{522, color.NRGBA{0x2D, 0xFF, 0x00, 0xFF}},
		{523, color.NRGBA{0x30, 0xFF, 0x00, 0xFF}},
		{524, color.NRGBA{0x34, 0xFF, 0x00, 0xFF}},
		{525, color.NRGBA{0x38, 0xFF, 0x00, 0xFF}},
		{526, color.NRGBA{0x3B, 0xFF, 0x00, 0xFF}},
		{527, color.NRGBA{0x3F, 0xFF, 0x00, 0xFF}},
		{528, color.NRGBA{0x43, 0xFF, 0x00, 0xFF}},
		{529, color.NRGBA{0x46, 0xFF, 0x00, 0xFF}},
		{530, color.NRGBA{0x4A, 0xFF, 0x00, 0xFF}},
		{531, color.NRGBA{0x4D, 0xFF, 0x00, 0xFF}},
		{532, color.NRGBA{0x51, 0xFF, 0x00, 0xFF}},
		{533, color.NRGBA{0x55, 0xFF, 0x00, 0xFF}},
		{534, color.NRGBA{0x58, 0xFF, 0x00, 0xFF}},
		{535, color.NRGBA{0x5C, 0xFF, 0x00, 0xFF}},
		{536, color.NRGBA{0x60, 0xFF, 0x00, 0xFF}},
		{537, color.NRGBA{0x63, 0xFF, 0x00, 0xFF}},
		{538, color.NRGBA{0x67, 0xFF, 0x00, 0xFF}},
		{539, color.NRGBA{0x6B, 0xFF, 0x00, 0xFF}},
		{540, color.NRGBA{0x6E, 0xFF, 0x00, 0xFF}},
		{541, color.NRGBA{0x72, 0xFF, 0x00, 0xFF}},
		{542, color.NRGBA{0x76, 0xFF, 0x00, 0xFF}},
		{543, color.NRGBA{0x79, 0xFF, 0x00, 0xFF}},
		{544, color.NRGBA{0x7D, 0xFF, 0x00, 0xFF}},
		{545, color.NRGBA{0x80, 0xFF, 0x00, 0xFF}},
		{546, color.NRGBA{0x84, 0xFF, 0x00, 0xFF}},
		{547, color.NRGBA{0x88, 0xFF, 0x00, 0xFF}},
		{548, color.NRGBA{0x8B, 0xFF, 0x00, 0xFF}},
		{549, color.NRGBA{0x8F, 0xFF, 0x00, 0xFF}},
		{550, color.NRGBA{0x93, 0xFF, 0x00, 0xFF}},
		{551, color.NRGBA{0x96, 0xFF, 0x00, 0xFF}},
		{552, color.NRGBA{0x9A, 0xFF, 0x00, 0xFF}},
		{553, color.NRGBA{0x9E, 0xFF, 0x00, 0xFF}},
		{554, color.NRGBA{0xA1, 0xFF, 0x00, 0xFF}},
		{555, color.NRGBA{0xA5, 0xFF, 0x00, 0xFF}},
		{556, color.NRGBA{0xA9, 0xFF, 0x00, 0xFF}},
		{557, color.NRGBA{0xAC, 0xFF, 0x00, 0xFF}},
		{558, color.NRGBA{0xB0, 0xFF, 0x00, 0xFF}},
		{559, color.NRGBA{0xB3, 0xFF, 0x00, 0xFF}},
		{560, color.NRGBA{0xB7, 0xFF, 0x00, 0xFF}},
		{561, color.NRGBA{0xBB, 0xFF, 0x00, 0xFF}},
		{562, color.NRGBA{0xBE, 0xFF, 0x00, 0xFF}},
		{563, color.NRGBA{0xC2, 0xFF, 0x00, 0xFF}},
		{564, color.NRGBA{0xC6, 0xFF, 0x00, 0xFF}},
		{565, color.NRGBA{0xC9, 0xFF, 0x00, 0xFF}},
		{566, color.NRGBA{0xCD, 0xFF, 0x00, 0xFF}},
		{567, color.NRGBA{0xD1, 0xFF, 0x00, 0xFF}},
		{568, color.NRGBA{0xD4, 0xFF, 0x00, 0xFF}},
		{569, color.NRGBA{0xD8, 0xFF, 0x00, 0xFF}},
		{570, color.NRGBA{0xDC, 0xFF, 0x00, 0xFF}},
		{571, color.NRGBA{0xDF, 0xFF, 0x00, 0xFF}},
		{572, color.NRGBA{0xE3, 0xFF, 0x00, 0xFF}},
		{573, color.NRGBA{0xE6, 0xFF, 0x00, 0xFF}},
		{574, color.NRGBA{0xEA, 0xFF, 0x00, 0xFF}},
		{575, color.NRGBA{0xEE, 0xFF, 0x00, 0xFF}},
		{576, color.NRGBA{0xF1, 0xFF, 0x00, 0xFF}},
		{577, color.NRGBA{0xF5, 0xFF, 0x00, 0xFF}},
		{578, color.NRGBA{0xF9, 0xFF, 0x00, 0xFF}},
		{579, color.NRGBA{0xFC, 0xFF, 0x00, 0xFF}},
		{580, color.NRGBA{0xFF, 0xFF, 0x00, 0xFF}},
		{581, color.NRGBA{0xFF, 0xFC, 0x00, 0xFF}},
		{582, color.NRGBA{0xFF, 0xF8, 0x00, 0xFF}},
		{583, color.NRGBA{0xFF, 0xF4, 0x00, 0xFF}},
		{584, color.NRGBA{0xFF, 0xF0, 0x00, 0xFF}},
		{585, color.NRGBA{0xFF, 0xEC, 0x00, 0xFF}},
		{586, color.NRGBA{0xFF, 0xE8, 0x00, 0xFF}},
		{587, color.NRGBA{0xFF, 0xE5, 0x00, 0xFF}},
		{588, color.NRGBA{0xFF, 0xE1, 0x00, 0xFF}},
		{589, color.NRGBA{0xFF, 0xDD, 0x00, 0xFF}},
		{590, color.NRGBA{0xFF, 0xD9, 0x00, 0xFF}},
		{591, color.NRGBA{0xFF, 0xD5, 0x00, 0xFF}},
		{592, color.NRGBA{0xFF, 0xD1, 0x00, 0xFF}},
		{593, color.NRGBA{0xFF, 0xCD, 0x00, 0xFF}},
		{594, color.NRGBA{0xFF, 0xC9, 0x00, 0xFF}},
		{595, color.NRGBA{0xFF, 0xC5, 0x00, 0xFF}},
		{596, color.NRGBA{0xFF, 0xC1, 0x00, 0xFF}},
		{597, color.NRGBA{0xFF, 0xBD, 0x00, 0xFF}},
		{598, color.NRGBA{0xFF, 0xB9, 0x00, 0xFF}},
		{599, color.NRGBA{0xFF, 0xB5, 0x00, 0xFF}},
		{600, color.NRGBA{0xFF, 0xB2, 0x00, 0xFF}},
		{601, color.NRGBA{0xFF, 0xAE, 0x00, 0xFF}},
		{602, color.NRGBA{0xFF, 0xAA, 0x00, 0xFF}},
		{603, color.NRGBA{0xFF, 0xA6, 0x00, 0xFF}},
		{604, color.NRGBA{0xFF, 0xA2, 0x00, 0xFF}},
		{605, color.NRGBA{0xFF, 0x9E, 0x00, 0xFF}},
		{606, color.NRGBA{0xFF, 0x9A, 0x00, 0xFF}},
		{607, color.NRGBA{0xFF, 0x96, 0x00, 0xFF}},
		{608, color.NRGBA{0xFF, 0x92, 0x00, 0xFF}},
		{609, color.NRGBA{0xFF, 0x8E, 0x00, 0xFF}},
		{610, color.NRGBA{0xFF, 0x8A, 0x00, 0xFF}},
		{611, color.NRGBA{0xFF, 0x86, 0x00, 0xFF}},
		{612, color.NRGBA{0xFF, 0x82, 0x00, 0xFF}},
		{613, color.NRGBA{0xFF, 0x7F, 0x00, 0xFF}},
		{614, color.NRGBA{0xFF, 0x7B, 0x00, 0xFF}},
		{615, color.NRGBA{0xFF, 0x77, 0x00, 0xFF}},
		{616, color.NRGBA{0xFF, 0x73, 0x00, 0xFF}},
		{617, color.NRGBA{0xFF, 0x6F, 0x00, 0xFF}},
		{618, color.NRGBA{0xFF, 0x6B, 0x00, 0xFF}},
		{619, color.NRGBA{0xFF, 0x67, 0x00, 0xFF}},
		{620, color.NRGBA{0xFF, 0x63, 0x00, 0xFF}},
		{621, color.NRGBA{0xFF, 0x5F, 0x00, 0xFF}},
		{622, color.NRGBA{0xFF, 0x5B, 0x00, 0xFF}},
		{623, color.NRGBA{0xFF, 0x57, 0x00, 0xFF}},
		{624, color.NRGBA{0xFF, 0x53, 0x00, 0xFF}},
		{625, color.NRGBA{0xFF, 0x4F, 0x00, 0xFF}},
		{626, color.NRGBA{0xFF, 0x4C, 0x00, 0xFF}},
		{627, color.NRGBA{0xFF, 0x48, 0x00, 0xFF}},
		{628, color.NRGBA{0xFF, 0x44, 0x00, 0xFF}},
		{629, color.NRGBA{0xFF, 0x40, 0x00, 0xFF}},
		{630, color.NRGBA{0xFF, 0x3C, 0x00, 0xFF}},
		{631, color.NRGBA{0xFF, 0x38, 0x00, 0xFF}},
		{632, color.NRGBA{0xFF, 0x34, 0x00, 0xFF}},
		{633, color.NRGBA{0xFF, 0x30, 0x00, 0xFF}},
		{634, color.NRGBA{0xFF, 0x2C, 0x00, 0xFF}},
		{635, color.NRGBA{0xFF, 0x28, 0x00, 0xFF}},
		{636, color.NRGBA{0xFF, 0x24, 0x00, 0xFF}},
		{637, color.NRGBA{0xFF, 0x20, 0x00, 0xFF}},
		{638, color.NRGBA{0xFF, 0x1C, 0x00, 0xFF}},
		{639, color.NRGBA{0xFF, 0x19, 0x00, 0xFF}},
		{640, color.NRGBA{0xFF, 0x15, 0x00, 0xFF}},
		{641, color.NRGBA{0xFF, 0x11, 0x00, 0xFF}},
		{642, color.NRGBA{0xFF, 0x0D, 0x00, 0xFF}},
		{643, color.NRGBA{0xFF, 0x09, 0x00, 0xFF}},
		{644, color.NRGBA{0xFF, 0x05, 0x00, 0xFF}},
		{645, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{646, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{647, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{648, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{649, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{650, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{651, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{652, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{653, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{654, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{655, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{656, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{657, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{658, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{659, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{660, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{661, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{662, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{663, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{664, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{665, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{666, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{667, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{668, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{669, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{670, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{671, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{672, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{673, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{674, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{675, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{676, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{677, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{678, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{679, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{680, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{681, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{682, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{683, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{684, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{685, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{686, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{687, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{688, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{689, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{690, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{691, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{692, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{693, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{694, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{695, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{696, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{697, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{698, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{699, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{700, color.NRGBA{0xFF, 0x00, 0x00, 0xFF}},
		{701, color.NRGBA{0xFF, 0x00, 0x00, 0xFD}},
		{702, color.NRGBA{0xFF, 0x00, 0x00, 0xFA}},
		{703, color.NRGBA{0xFF, 0x00, 0x00, 0xF7}},
		{704, color.NRGBA{0xFF, 0x00, 0x00, 0xF5}},
		{705, color.NRGBA{0xFF, 0x00, 0x00, 0xF2}},
		{706, color.NRGBA{0xFF, 0x00, 0x00, 0xEF}},
		{707, color.NRGBA{0xFF, 0x00, 0x00, 0xEC}},
		{708, color.NRGBA{0xFF, 0x00, 0x00, 0xE9}},
		{709, color.NRGBA{0xFF, 0x00, 0x00, 0xE6}},
		{710, color.NRGBA{0xFF, 0x00, 0x00, 0xE3}},
		{711, color.NRGBA{0xFF, 0x00, 0x00, 0xE0}},
		{712, color.NRGBA{0xFF, 0x00, 0x00, 0xDE}},
		{713, color.NRGBA{0xFF, 0x00, 0x00, 0xDB}},
		{714, color.NRGBA{0xFF, 0x00, 0x00, 0xD8}},
		{715, color.NRGBA{0xFF, 0x00, 0x00, 0xD5}},
		{716, color.NRGBA{0xFF, 0x00, 0x00, 0xD2}},
		{717, color.NRGBA{0xFF, 0x00, 0x00, 0xCF}},
		{718, color.NRGBA{0xFF, 0x00, 0x00, 0xCC}},
		{719, color.NRGBA{0xFF, 0x00, 0x00, 0xC9}},
		{720, color.NRGBA{0xFF, 0x00, 0x00, 0xC7}},
		{721, color.NRGBA{0xFF, 0x00, 0x00, 0xC4}},
		{722, color.NRGBA{0xFF, 0x00, 0x00, 0xC1}},
		{723, color.NRGBA{0xFF, 0x00, 0x00, 0xBE}},
		{724, color.NRGBA{0xFF, 0x00, 0x00, 0xBB}},
		{725, color.NRGBA{0xFF, 0x00, 0x00, 0xB8}},
		{726, color.NRGBA{0xFF, 0x00, 0x00, 0xB5}},
		{727, color.NRGBA{0xFF, 0x00, 0x00, 0xB3}},
		{728, color.NRGBA{0xFF, 0x00, 0x00, 0xB0}},
		{729, color.NRGBA{0xFF, 0x00, 0x00, 0xAD}},
		{730, color.NRGBA{0xFF, 0x00, 0x00, 0xAA}},
		{731, color.NRGBA{0xFF, 0x00, 0x00, 0xA7}},
		{732, color.NRGBA{0xFF, 0x00, 0x00, 0xA4}},
		{733, color.NRGBA{0xFF, 0x00, 0x00, 0xA1}},
		{734, color.NRGBA{0xFF, 0x00, 0x00, 0x9E}},
		{735, color.NRGBA{0xFF, 0x00, 0x00, 0x9C}},
		{736, color.NRGBA{0xFF, 0x00, 0x00, 0x99}},
		{737, color.NRGBA{0xFF, 0x00, 0x00, 0x96}},
		{738, color.NRGBA{0xFF, 0x00, 0x00, 0x93}},
		{739, color.NRGBA{0xFF, 0x00, 0x00, 0x90}},
		{740, color.NRGBA{0xFF, 0x00, 0x00, 0x8D}},
		{741, color.NRGBA{0xFF, 0x00, 0x00, 0x8A}},
		{742, color.NRGBA{0xFF, 0x00, 0x00, 0x88}},
		{743, color.NRGBA{0xFF, 0x00, 0x00, 0x85}},
		{744, color.NRGBA{0xFF, 0x00, 0x00, 0x82}},
		{745, color.NRGBA{0xFF, 0x00, 0x00, 0x7F}},
		{746, color.NRGBA{0xFF, 0x00, 0x00, 0x7C}},
		{747, color.NRGBA{0xFF, 0x00, 0x00, 0x79}},
		{748, color.NRGBA{0xFF, 0x00, 0x00, 0x76}},
		{749, color.NRGBA{0xFF, 0x00, 0x00, 0x73}},
		{750, color.NRGBA{0xFF, 0x00, 0x00, 0x71}},
		{751, color.NRGBA{0xFF, 0x00, 0x00, 0x6E}},
		{752, color.NRGBA{0xFF, 0x00, 0x00, 0x6B}},
		{753, color.NRGBA{0xFF, 0x00, 0x00, 0x68}},
		{754, color.NRGBA{0xFF, 0x00, 0x00, 0x65}},
		{755, color.NRGBA{0xFF, 0x00, 0x00, 0x62}},
		{756, color.NRGBA{0xFF, 0x00, 0x00, 0x5F}},
		{757, color.NRGBA{0xFF, 0x00, 0x00, 0x5C}},
		{758, color.NRGBA{0xFF, 0x00, 0x00, 0x5A}},
		{759, color.NRGBA{0xFF, 0x00, 0x00, 0x57}},
		{760, color.NRGBA{0xFF, 0x00, 0x00, 0x54}},
		{761, color.NRGBA{0xFF, 0x00, 0x00, 0x51}},
		{762, color.NRGBA{0xFF, 0x00, 0x00, 0x4E}},
		{763, color.NRGBA{0xFF, 0x00, 0x00, 0x4B}},
		{764, color.NRGBA{0xFF, 0x00, 0x00, 0x48}},
		{765, color.NRGBA{0xFF, 0x00, 0x00, 0x46}},
		{766, color.NRGBA{0xFF, 0x00, 0x00, 0x43}},
		{767, color.NRGBA{0xFF, 0x00, 0x00, 0x40}},
		{768, color.NRGBA{0xFF, 0x00, 0x00, 0x3D}},
		{769, color.NRGBA{0xFF, 0x00, 0x00, 0x3A}},
		{770, color.NRGBA{0xFF, 0x00, 0x00, 0x37}},
		{771, color.NRGBA{0xFF, 0x00, 0x00, 0x34}},
		{772, color.NRGBA{0xFF, 0x00, 0x00, 0x31}},
		{773, color.NRGBA{0xFF, 0x00, 0x00, 0x2F}},
		{774, color.NRGBA{0xFF, 0x00, 0x00, 0x2C}},
		{775, color.NRGBA{0xFF, 0x00, 0x00, 0x29}},
		{776, color.NRGBA{0xFF, 0x00, 0x00, 0x26}},
		{777, color.NRGBA{0xFF, 0x00, 0x00, 0x23}},
		{778, color.NRGBA{0xFF, 0x00, 0x00, 0x20}},
		{779, color.NRGBA{0xFF, 0x00, 0x00, 0x1D}},
		{780, color.NRGBA{0xFF, 0x00, 0x00, 0x1A}},
		{781, color.NRGBA{0x00, 0x00, 0x00, 0x00}},
	}
	/*
		for _, line := range data {
			c := waveLength2RGB(line.input)
			fmt.Printf("{%d, color.NRGBA{0x%02X, 0x%02X, 0x%02X, 0x%02X}},\n", int(line.input), c.R, c.G, c.B, c.A)
		}
	*/
	for i, line := range data {
		ut.AssertEqualIndex(t, i, line.expected, waveLength2RGB(line.input))
	}
}

//

type expectation struct {
	offset time.Duration
	colors []color.NRGBA
}

func frames(t *testing.T, p Pattern, expectations []expectation) {
	pixels := make([]color.NRGBA, len(expectations[0].colors))
	for frame, e := range expectations {
		p.NextFrame(pixels, e.offset)
		for j := range e.colors {
			a := e.colors[j]
			b := pixels[j]
			dR := int(a.R) - int(b.R)
			dG := int(a.G) - int(b.G)
			dB := int(a.B) - int(b.B)
			dA := int(a.A) - int(b.A)
			if dR > 1 || dR < -1 || dG > 1 || dG < -1 || dB > 1 || dB < -1 || dA > 1 || dA < -1 {
				t.Fatalf("frame=%d; pixel=%d; %v != %v", frame, j, a, b)
			}
		}
	}
}
