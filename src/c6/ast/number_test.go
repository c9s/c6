package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNumberAddInt(t *testing.T) {
	var num = NewNumber(100, nil)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, 200.0, num.Value)
}

func TestNumberPxUnitString(t *testing.T) {
	var num = NewNumber(100, nil)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, "200px", num.String())
}

func TestNumberPtUnitString(t *testing.T) {
	var num = NewNumber(100, nil)
	num.AddInt(100)
	num.SetUnit(UNIT_PT)
	assert.Equal(t, "200pt", num.String())
}

func TestNumberEmUnitString(t *testing.T) {
	var num = NewNumber(100, nil)
	num.AddInt(100)
	num.SetUnit(UNIT_EM)
	assert.Equal(t, "200em", num.String())
}

func TestNumberAddIntInt(t *testing.T) {
	var a = NewNumber(10, nil)
	var b = NewNumber(23, nil)
	var c = NumberAdd(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33.0, c.Value)
}

func TestNumberAddIntFloat(t *testing.T) {
	var a = NewNumber(10, nil)
	var b = NewNumber(23.3, nil)
	var c = NumberAdd(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33.3, c.Value)
}

func TestNumberAddFloatFloat(t *testing.T) {
	var a = NewNumber(10.3, nil)
	var b = NewNumber(23.3, nil)
	var c = NumberAdd(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 33.6, c.Value)
}

func TestNumberAddFloatInt(t *testing.T) {
	var a = NewNumber(10.3, nil)
	var b = NewNumber(3, nil)
	var c = NumberAdd(a, b)
	assert.Equal(t, 13.3, c.Value)
}

func TestNumberAddFloatPxFloatPx(t *testing.T) {
	var a = NewNumber(10.3, nil)
	a.SetUnit(UNIT_PX)
	var b = NewNumber(3, nil)
	a.SetUnit(UNIT_PX)

	var c = NumberAdd(a, b)
	assert.NotNil(t, c)
	assert.Equal(t, 13.3, c.Value)
	assert.Equal(t, "13.3px", c.String())
}
