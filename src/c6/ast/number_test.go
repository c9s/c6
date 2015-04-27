package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNumberAddInt(t *testing.T) {
	var num = NewIntegerNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, 200, num.Int)
}

func TestNumberPxUnitString(t *testing.T) {
	var num = NewIntegerNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	assert.Equal(t, "200px", num.String())
}

func TestNumberPtUnitString(t *testing.T) {
	var num = NewIntegerNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PT)
	assert.Equal(t, "200pt", num.String())
}

func TestNumberEmUnitString(t *testing.T) {
	var num = NewIntegerNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_EM)
	assert.Equal(t, "200em", num.String())
}

func TestNumberAddIntInt(t *testing.T) {
	var a = NewIntegerNumber(10)
	var b = NewIntegerNumber(23)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	intNum, ok := c.(*IntegerNumber)
	assert.True(t, ok)
	assert.Equal(t, 33, intNum.Int)
}

func TestNumberAddIntFloat(t *testing.T) {
	var a = NewIntegerNumber(10)
	var b = NewFloatNumber(23.3)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	num, ok := c.(*FloatNumber)
	assert.True(t, ok)
	assert.Equal(t, 33.3, num.Float)
}

func TestNumberAddFloatFloat(t *testing.T) {
	var a = NewFloatNumber(10.3)
	var b = NewFloatNumber(23.3)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	num, ok := c.(*FloatNumber)
	assert.True(t, ok)
	assert.Equal(t, 33.6, num.Float)
}

func TestNumberAddFloatInt(t *testing.T) {
	var a = NewFloatNumber(10.3)
	var b = NewIntegerNumber(3)
	var c = AddNumber(a, b)
	assert.NotNil(t, c)
	num, ok := c.(*FloatNumber)
	assert.True(t, ok)
	assert.Equal(t, 13.3, num.Float)
}
