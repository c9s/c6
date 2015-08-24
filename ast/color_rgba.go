package ast

import "fmt"

type RGBAColor struct {
	R     uint32
	G     uint32
	B     uint32
	A     float32
	Token *Token
}

// Factor functions
func NewRGBAColorWithHexCode(hex string, token *Token) *RGBAColor {
	var r, g, b, a = HexToRGBA(hex)
	return &RGBAColor{r, g, b, a, token}
}

func (self RGBAColor) CanBeNode()  {}
func (self RGBAColor) CanBeColor() {}

func (self RGBAColor) Boolean() bool {
	return true
}

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

/*
RGBColor can present rgb(....)

*/
type RGBColor struct {
	R     uint32
	G     uint32
	B     uint32
	Token *Token
}

func (self RGBColor) CanBeNode()  {}
func (self RGBColor) CanBeColor() {}

func (self RGBColor) Boolean() bool {
	return true
}

func (self RGBColor) Hex() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func (self RGBColor) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", self.R, self.G, self.B)
}

func NewRGBColor(r, g, b uint32, token *Token) *RGBColor {
	return &RGBColor{r, g, b, token}
}

func NewRGBColorWithHexCode(hex string, token *Token) *RGBColor {
	var r, g, b, _ = HexToRGBA(hex)
	return &RGBColor{r, g, b, token}
}
