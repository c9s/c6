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

// hsv() is not supported in CSS3, we need to convert it to hex color
func (self HSVColor) String() string {
	return fmt.Sprintf("hsv(%G, %G, %G)", self.H, self.S, self.V)
}

func NewHSVColor(h, s, v float64, token *Token) *HSVColor {
	return &HSVColor{h, s, v, token}
}

func RGBToHSV(ir, ig, ib uint) (h, s, v float64) {
	// cast to float64 for math.* API
	var r = float64(ir)
	var g = float64(ig)
	var b = float64(ib)

	var min = math.Min(math.Min(r, g), b)

	v = math.Max(math.Max(r, g), b)
	var C = v - min

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
