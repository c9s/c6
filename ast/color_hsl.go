package ast

/*
http://stackoverflow.com/questions/2353211/hsl-to-rgb-color-conversion
*/

import "fmt"
import "math"

type HSLColor struct {
	H     float64
	S     float64
	L     float64
	Token *Token
}

func (self HSLColor) CanBeColor() {}
func (self HSLColor) CanBeNode()  {}
func (self HSLColor) HSLAColor() *HSLAColor {
	return NewHSLAColor(self.H, self.S, self.L, 0, nil)

}
func (c HSLColor) RGBAColor() *RGBAColor {
	var r, g, b = HSLToRGB(c.H, c.S, c.L)
	return NewRGBAColor(uint32(r), uint32(g), uint32(b), 1, nil)
}
func (c HSLColor) RGBColor() *RGBColor {
	var r, g, b = HSLToRGB(c.H, c.S, c.L)
	return NewRGBColor(uint32(r), uint32(g), uint32(b), nil)
}

func (self HSLColor) String() string {
	return fmt.Sprintf("hsl(%G, %G, %G)", self.H, self.S, self.L)
}

func NewHSLColor(h, s, v float64, token *Token) *HSLColor {
	return &HSLColor{h, s, v, token}
}

type HSLAColor struct {
	H     float64
	S     float64
	L     float64
	A     float64
	Token *Token
}

func (self HSLAColor) CanBeColor() {}
func (self HSLAColor) CanBeNode()  {}
func (self HSLAColor) String() string {
	return fmt.Sprintf("hsl(%G, %G, %G, %G)", self.H, self.S, self.L, self.A)
}

func (self HSLAColor) Boolean() bool {
	return true
}

func NewHSLAColor(h, s, v, a float64, token *Token) *HSLAColor {
	return &HSLAColor{h, s, v, a, token}
}

/*
h = 0~360
s = 0~1
l = 0~1
*/
func HSLToRGB(h, s, l float64) (r, g, b uint32) {
	if h > 0 {
		h = h / 360
	}
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		fR = ConvertHUE(p, q, h+1.0/3)
		fG = ConvertHUE(p, q, h)
		fB = ConvertHUE(p, q, h-1.0/3)
	}
	r = uint32((fR * 255) + 0.5)
	g = uint32((fG * 255) + 0.5)
	b = uint32((fB * 255) + 0.5)
	return
}

func ConvertHUE(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}

func RGBToHSL(r, g, b uint32) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	l = (max + min) / 2
	if max == min {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}
