package ast

import "strconv"
import "fmt"

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

func (self HexColor) CanBeNode()       {}
func (self HexColor) CanBeColor()      {}
func (self HexColor) String() string {
	return "#" + string(self.Hex)
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
			fmt.Printf("%+v", rgb)
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
	R     uint32
	G     uint32
	B     uint32
	A     float32
	Token *Token
}

func (self RGBAColor) CanBeNode()  {}
func (self RGBAColor) CanBeColor() {}

// NOTE: 8 char hex color is only supported by IE.
func (self RGBAColor) Hex() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func (self RGBAColor) String() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %g)", self.R, self.G, self.B, self.A)
}

func NewRGBAColor(r, g, b uint32, a float32, token *Token) *RGBAColor {
	return &RGBAColor{r, g, b, a, token}
}

type RGBColor struct {
	R     uint32
	G     uint32
	B     uint32
	Token *Token
}

func (self RGBColor) CanBeNode()       {}
func (self RGBColor) CanBeColor()      {}


func (self RGBColor) Hex() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func (self RGBColor) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", self.R, self.G, self.B)
}

func NewRGBColor(r, g, b uint32, token *Token) *RGBColor {
	return &RGBColor{r, g, b, token}
}
