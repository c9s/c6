package ast

import "strconv"

type Hex string

func (hex Hex) String() string {
	return string(hex)
}

type Color interface {
	CanBeColor()
}

type HexColor struct {
	Hex   Hex
	R     uint32
	G     uint32
	B     uint32
	Token *Token
}

func (self HexColor) CanBeNode() {}

func (self HexColor) CanBeColor() {}

func (self HexColor) Boolean() bool {
	return true
}

func (self HexColor) String() string {
	return "#" + string(self.Hex)
}

func NewHexColorFromToken(token *Token) *HexColor {
	return NewHexColor(token.Str, token)
}

func NewHexColor(hex string, token *Token) *HexColor {
	var r, g, b, _ = HexToRGBA(hex)
	return &HexColor{Hex(hex), r, g, b, token}
}

// HexToRGB converts an Hex string to a RGB triple.
func HexToRGBA(h string) (uint32, uint32, uint32, float32) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) == 3 {
		// rebuild hex string
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	if len(h) == 6 {
		if rgb, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint32(rgb >> 16), uint32((rgb >> 8) & 0xFF), uint32(rgb & 0xFF), 0
		}
	}
	if len(h) == 8 {
		if rgba, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint32(rgba >> 24), uint32(rgba >> 16 & 0xFF), uint32((rgba >> 8) & 0xFF), float32(rgba&0xFF) / 255
		}
	}
	return 0, 0, 0, 0
}
