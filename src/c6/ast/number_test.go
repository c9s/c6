package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNumberAddInt(t *testing.T) {
	var num = NewNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, 200, num.Int)
}

func TestNumberPxUnitString(t *testing.T) {
	var num = NewNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, "200px", num.String())
}

func TestNumberPtUnitString(t *testing.T) {
	var num = NewNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PT)
	assert.Equal(t, "200pt", num.String())
}

func TestNumberEmUnitString(t *testing.T) {
	var num = NewNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_EM)
	assert.Equal(t, "200em", num.String())
}

func TestNumberAddIntInt(t *testing.T) {
	var a = NewNumber(10)
	var b = NewNumber(23)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33, c.Int)
}

func TestNumberAddIntFloat(t *testing.T) {
	var a = NewNumber(10)
	var b = NewNumber(23.3)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33.3, c.Value)
}

func TestNumberAddFloatFloat(t *testing.T) {
	var a = NewNumber(10.3)
	var b = NewNumber(23.3)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33.6, c.Float)
}

func TestNumberAddFloatInt(t *testing.T) {
	var a = NewFloatNumber(10.3)
	var b = NewIntegerNumber(3)
	var c = AddNumber(a, b)
	assert.Equal(t, 13.3, c.Float)
}

func TestNumberAddFloatPxFloatPx(t *testing.T) {
	var a = NewNumber(10.3)
	a.SetUnit(UNIT_PX)
	var b = NewNumber(3)
	a.SetUnit(UNIT_PX)

	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 13.3, c.Float)
	assert.Equal(t, "13.3px", c.String())
}
