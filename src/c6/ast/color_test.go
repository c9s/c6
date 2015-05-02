package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewHexColor(t *testing.T) {
	var c = NewHexColor("#ffffff", nil)

	assert.Equal(t, uint8(255), c.R)
	assert.Equal(t, uint8(255), c.G)
	assert.Equal(t, uint8(255), c.B)
	assert.Equal(t, c.Hex, Hex("#ffffff"))
}

func TestHex6CharToRGBA(t *testing.T) {
	var r, g, b, a = HexToRGBA("#FF2030")
	assert.Equal(t, uint8(255), r)
	assert.Equal(t, uint8(32), g)
	assert.Equal(t, uint8(48), b)
	assert.Equal(t, float32(0.0), a)
}

func TestHex8CharToRGBA(t *testing.T) {
	var r, g, b, a = HexToRGBA("#FF203040")
	assert.Equal(t, uint8(255), r)
	assert.Equal(t, uint8(32), g)
	assert.Equal(t, uint8(48), b)
	// float32 is not precise, hence use NotEqual here.
	assert.NotEqual(t, float32(0.0), a)
}
