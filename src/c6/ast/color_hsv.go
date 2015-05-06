package ast

import "fmt"
import "math"

type HSVColor struct {
	H     float64
	S     float64
	V     float64
	Token *Token
}

func (self HSVColor) CanBeColor() {}
func (self HSVColor) CanBeNode()  {}

func (self HSVColor) Boolean() bool {
	return true
}

func (c HSVColor) RGBAColor() *RGBAColor {
	r, g, b := HSVToRGB(c.H, c.S, c.V)
	return NewRGBAColor(uint32(r), uint32(g), uint32(b), 1, nil)
}

func (c HSVColor) RGBColor() *RGBColor {
	r, g, b := HSVToRGB(c.H, c.S, c.V)
	return NewRGBColor(uint32(r), uint32(g), uint32(b), nil)
}

// hsv() is not supported in CSS3, we need to convert it to hex color
func (self HSVColor) String() string {
	return fmt.Sprintf("hsv(%G, %G, %G)", self.H, self.S, self.V)
}

func NewHSVColor(h, s, v float64, token *Token) *HSVColor {
	return &HSVColor{h, s, v, token}
}

func RGBToHSV(ir, ig, ib uint32) (h, s, v float64) {
	r := float64(ir) / 255
	g := float64(ig) / 255
	b := float64(ib) / 255

	min := math.Min(math.Min(r, g), b)
	v = math.Max(math.Max(r, g), b)
	C := v - min

	s = 0.0
	if v != 0.0 {
		s = C / v
	}

	h = 0.0 // We use 0 instead of undefined as in wp.
	if min != v {
		if v == r {
			h = math.Mod((g-b)/C, 6.0)
		}
		if v == g {
			h = (b-r)/C + 2.0
		}
		if v == b {
			h = (r-g)/C + 4.0
		}
		h *= 60.0
		if h < 0.0 {
			h += 360.0
		}
	}
	return
}

func HSVToRGB(h, s, v float64) (r, g, b uint32) {
	var fR, fG, fB float64
	i := math.Floor(h / 60)
	f := h/60 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)
	switch int(i) % 6 {
	case 0:
		fR, fG, fB = v, t, p
	case 1:
		fR, fG, fB = q, v, p
	case 2:
		fR, fG, fB = p, v, t
	case 3:
		fR, fG, fB = p, q, v
	case 4:
		fR, fG, fB = t, p, v
	case 5:
		fR, fG, fB = v, p, q
	}
	fmt.Printf("%G %G %G\n", fR, fG, fB)
	r = uint32((fR * 255) + 0.5)
	g = uint32((fG * 255) + 0.5)
	b = uint32((fB * 255) + 0.5)
	return
}
