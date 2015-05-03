package ast

import "fmt"
import "math"

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

// NOTE: 8 char hex color is only supported by IE.
func (self RGBAColor) Hex() Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", self.R, self.G, self.B))
}

func (self RGBAColor) String() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %g)", self.R, self.G, self.B, self.A)
}

func RGBAColorAddNumber(c *RGBAColor, n *Number) *RGBAColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorSubNumber(c *RGBAColor, n *Number) *RGBAColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorMulNumber(c *RGBAColor, n *Number) *RGBAColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorDivNumber(c *RGBAColor, n *Number) *RGBAColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return NewRGBAColor(uint32(r), uint32(g), uint32(b), c.A, nil)
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

func RGBColorAddNumber(c *RGBColor, n *Number) *RGBColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return NewRGBColor(r, g, b, nil)
}

func RGBColorSubNumber(c *RGBColor, n *Number) *RGBColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return NewRGBColor(r, g, b, nil)
}

func RGBColorMulNumber(c *RGBColor, n *Number) *RGBColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return NewRGBColor(r, g, b, nil)
}

func RGBColorDivNumber(c *RGBColor, n *Number) *RGBColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return NewRGBColor(uint32(r), uint32(g), uint32(b), nil)
}
