package ast

import "fmt"

type Number interface {
	SetUnit(unit UnitType)
	IsNumber()
	CanBeExpression()
}

type FloatNumber struct {
	Float float64
	Unit  UnitType
	Token *Token
}

func (num *FloatNumber) AddFloat(a float64) {
	num.Float += a
}

func (num *FloatNumber) AddInt(a int) {
	num.Float += float64(a)
}

func NewFloatNumber(num float64) *FloatNumber {
	return &FloatNumber{num, UNIT_NONE, nil}
}

func (num FloatNumber) IsNumber()        {}
func (num FloatNumber) CanBeExpression() {}

func (num FloatNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num FloatNumber) String() (out string) {
	out += fmt.Sprintf("%.2f", num.Float)
	if num.Unit != UNIT_NONE {
		out += num.Unit.String()
	}
	return out
}

type IntegerNumber struct {
	Int   int64
	Unit  UnitType
	Token *Token
}

func NewIntegerNumber(num int64) *IntegerNumber {
	return &IntegerNumber{num, UNIT_NONE, nil}
}

func (num IntegerNumber) IsNumber()        {}
func (num IntegerNumber) CanBeExpression() {}

func (num *IntegerNumber) AddFloat(a float64) {
	num.Int += int64(a)
}

func (num *IntegerNumber) AddInt(a int64) {
	num.Int += a
}

func (num IntegerNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num IntegerNumber) String() (out string) {
	out += fmt.Sprintf("%d", num.Int)
	if num.Unit > 0 {
		out += num.Unit.String()
	}
	return out
}
