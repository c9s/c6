package ast

import "fmt"

type IntegerNumber struct {
	Int   int64
	Unit  UnitType
	Token *Token
}

func NewIntegerNumber(num int64) *IntegerNumber {
	return &IntegerNumber{num, UNIT_NONE, nil}
}

func (num IntegerNumber) CanBeNumber()     {}
func (num IntegerNumber) CanBeNode()       {}
func (num IntegerNumber) CanBeExpression() {}

func (num *IntegerNumber) GetUnit() UnitType {
	return num.Unit
}

func (num *IntegerNumber) AddFloat(a float64) {
	num.Int += int64(a)
}

func (num *IntegerNumber) AddInt(a int64) {
	num.Int += a
}

func (num *IntegerNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num IntegerNumber) String() (out string) {
	out += fmt.Sprintf("%d", num.Int)
	if num.Unit > 0 {
		out += num.Unit.UnitString()
	}
	return out
}
