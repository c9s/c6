package ast

import "strconv"
import "fmt"

type Hex string

type Color interface {
	CanBeColor()
}

type HexColor struct {
	Hex   Hex
	R     uint8
	G     uint8
	B     uint8
	Token *Token
}

func (self HexColor) CanBeExpression() {}
func (self HexColor) CanBeNode()       {}
func (self HexColor) CanBeColor()      {}
func (self HexColor) CanBeValue()      {}
func (self HexColor) String() string {
	return "#" + string(self.Hex)
}

func NewHexColor(hex string, token *Token) *HexColor {
	var r, g, b, _ = HexToRGBA(hex)
	return &HexColor{Hex(hex), r, g, b, token}
}

// HexToRGB converts an Hex string to a RGB triple.
func HexToRGBA(h string) (uint8, uint8, uint8, float32) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) == 3 {
		// rebuild hex string
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	if len(h) == 6 {
		if rgb, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			fmt.Printf("%+v", rgb)
			return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF), 0
		}
	}
	if len(h) == 8 {
		if rgba, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint8(rgba >> 24), uint8(rgba >> 16), uint8((rgba >> 8) & 0xFF), float32(rgba&0xFF) / 255
		}
	}
	return 0, 0, 0, 0
}

func NewRGBColorWithHexCode(hex string, token *Token) *RGBColor {
	var r, g, b, _ = HexToRGBA(hex)
	return &RGBColor{r, g, b, token}
}

// Factor functions
func NewRGBAColorWithHexCode(hex string, token *Token) *RGBAColor {
	var r, g, b, a = HexToRGBA(hex)
	return &RGBAColor{r, g, b, a, token}
}

type RGBAColor struct {
	R     uint8
	G     uint8
	B     uint8
	A     float32
	Token *Token
}

func (self RGBAColor) CanBeExpression() {}
func (self RGBAColor) CanBeNode()       {}
func (self RGBAColor) CanBeColor()      {}
func (self RGBAColor) CanBeValue()      {}

// NOTE: 8 char hex color is only supported by IE.
func (self RGBAColor) ToHexColor() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func NewRGBAColor(r, g, b uint8, a float32, token *Token) *RGBAColor {
	return &RGBAColor{r, g, b, a, token}
}

type RGBColor struct {
	R     uint8
	G     uint8
	B     uint8
	Token *Token
}

func (self RGBColor) CanBeExpression() {}
func (self RGBColor) CanBeNode()       {}
func (self RGBColor) CanBeColor()      {}

func (self RGBColor) CanBeValue() {}

func (self RGBColor) ToHexColor() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func NewRGBColor(r, g, b uint8, token *Token) *RGBColor {
	return &RGBColor{r, g, b, token}
}
