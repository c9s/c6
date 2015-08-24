package runtime

import (
	"fmt"
	"math"

	"github.com/c9s/c6/ast"
)

func HexColorAddNumber(c *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := c.R + uint32(num.Value)
	g := c.G + uint32(num.Value)
	b := c.B + uint32(num.Value)
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{
		Hex:   hex,
		R:     r,
		G:     g,
		B:     b,
		Token: nil,
	}
}

func uintsub(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return 0
}

func HexColorSubNumber(c *ast.HexColor, num *ast.Number) *ast.HexColor {
	val := uint32(num.Value)
	r := uintsub(c.R, val)
	g := uintsub(c.G, val)
	b := uintsub(c.B, val)
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{
		Hex:   hex,
		R:     r,
		G:     g,
		B:     b,
		Token: nil,
	}
}

func HexColorMulNumber(color *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := uint32(math.Floor(float64(color.R) * num.Value))
	g := uint32(math.Floor(float64(color.G) * num.Value))
	b := uint32(math.Floor(float64(color.B) * num.Value))
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{
		Hex:   hex,
		R:     r,
		G:     g,
		B:     b,
		Token: nil,
	}
}

func HexColorDivNumber(color *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := uint32(math.Floor(float64(color.R) / num.Value))
	g := uint32(math.Floor(float64(color.G) / num.Value))
	b := uint32(math.Floor(float64(color.B) / num.Value))
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{
		Hex:   hex,
		R:     r,
		G:     g,
		B:     b,
		Token: nil,
	}
}

func RGBColorAddNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorSubNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorMulNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorDivNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return ast.NewRGBColor(uint32(r), uint32(g), uint32(b), nil)
}

func RGBAColorAddNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorSubNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorMulNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorDivNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return ast.NewRGBAColor(uint32(r), uint32(g), uint32(b), c.A, nil)
}
