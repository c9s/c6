package ast

import "strconv"

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

func (num *FloatNumber) GetUnit() UnitType {
	return num.Unit
}

func NewFloatNumber(num float64) *FloatNumber {
	return &FloatNumber{num, UNIT_NONE, nil}
}

func (num FloatNumber) CanBeNumber()     {}
func (num FloatNumber) CanBeNode()       {}
func (num FloatNumber) CanBeExpression() {}

func (num *FloatNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num FloatNumber) String() (out string) {
	out += strconv.FormatFloat(num.Float, 'G', -1, 64)
	if num.Unit != UNIT_NONE {
		out += num.Unit.UnitString()
	}
	return out
}
