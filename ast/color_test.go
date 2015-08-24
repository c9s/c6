package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewRGBAColor(t *testing.T) {
	var c = NewRGBAColor(10, 20, 30, 0.5, nil)
	assert.Equal(t, "rgba(10, 20, 30, 0.5)", c.String())
}

func TestNewRGBAColor2(t *testing.T) {
	var c = NewRGBAColor(10, 20, 30, 0.22, nil)
	assert.Equal(t, "rgba(10, 20, 30, 0.22)", c.String())
	assert.Equal(t, Hex("#0A141E"), c.Hex())
}

func TestNewRGBColor(t *testing.T) {
	var c = NewRGBColor(10, 20, 30, nil)
	assert.Equal(t, "rgb(10, 20, 30)", c.String())
	assert.Equal(t, Hex("#0A141E"), c.Hex())
}

func TestNewHexColor(t *testing.T) {
	var c = NewHexColor("#ffffff", nil)
	assert.Equal(t, uint32(255), c.R)
	assert.Equal(t, uint32(255), c.G)
	assert.Equal(t, uint32(255), c.B)
	assert.Equal(t, c.Hex, Hex("#ffffff"))
}

func TestHex6CharToRGBA(t *testing.T) {
	var r, g, b, a = HexToRGBA("#FF2030")
	assert.Equal(t, uint32(255), r)
	assert.Equal(t, uint32(32), g)
	assert.Equal(t, uint32(48), b)
	assert.Equal(t, float32(0.0), a)
}

func TestHSLColorRedToRGB(t *testing.T) {
	// hsl(0,  100%,50%) = red
	var c = NewHSLColor(0, 1, 0.5, nil)
	assert.NotNil(t, c)
	assert.Equal(t, "rgb(255, 0, 0)", c.RGBColor().String())
}

func TestHSLColorGreenToRGB(t *testing.T) {
	// hsl(120,100%,50%) = green
	var c = NewHSLColor(120, 1, 0.5, nil)
	assert.NotNil(t, c)
	assert.Equal(t, "rgb(0, 255, 0)", c.RGBColor().String())
}

func TestHSVColorRedToRGB(t *testing.T) {
	var c = NewHSVColor(0, 1, 1, nil)
	assert.Equal(t, "rgb(255, 0, 0)", c.RGBColor().String())
}

func TestHSVColorGreenToRGB(t *testing.T) {
	var c = NewHSVColor(120, 0.5, 1, nil)
	assert.Equal(t, "rgb(128, 255, 128)", c.RGBColor().String())

	var h, s, v = RGBToHSV(0, 255, 0)
	assert.Equal(t, 120.0, h)
	assert.Equal(t, 1.0, s)
	assert.Equal(t, 1.0, v)
}

func TestHex8CharToRGBA(t *testing.T) {
	var r, g, b, a = HexToRGBA("#FF203040")
	assert.Equal(t, uint32(255), r)
	assert.Equal(t, uint32(32), g)
	assert.Equal(t, uint32(48), b)
	// float32 is not precise, hence use NotEqual here.
	assert.NotEqual(t, float32(0.0), a)
}
