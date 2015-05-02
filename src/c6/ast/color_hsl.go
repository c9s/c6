package ast

import "fmt"

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

func NewHSLAColor(h, s, v, a float64, token *Token) *HSLAColor {
	return &HSLAColor{h, s, v, a, token}
}

func HSLToRGB(h, s, l float64) (r, g, b uint8) {
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
		fR = HUEToRGB(p, q, h+1.0/3)
		fG = HUEToRGB(p, q, h)
		fB = HUEToRGB(p, q, h-1.0/3)
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

func HUEToRGB(p, q, t float64) float64 {
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
