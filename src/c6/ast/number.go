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
}

func (num FloatNumber) AddFloat(a float64) {
	num.Float += a
}

func (num FloatNumber) AddInt(a int) {
	num.Float += float64(a)
}

func NewFloatNumber(num float64) *FloatNumber {
	return &FloatNumber{num, 0}
}

func (num FloatNumber) IsNumber()        {}
func (num FloatNumber) CanBeExpression() {}

func (num FloatNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num FloatNumber) String() (out string) {
	out += fmt.Sprintf("%.2f", num.Float)

	if num.Unit > 0 {
		switch num.Unit {
		case UNIT_PX:
			out += "px"
		case UNIT_PT:
			out += "pt"
		case UNIT_EM:
			out += "em"
		case UNIT_REM:
			out += "rem"
		default:
			panic("unimplemented type")
		}
	}
	return out
}

type IntegerNumber struct {
	Int  int64
	Unit UnitType
}

func NewIntegerNumber(num int64) *IntegerNumber {
	return &IntegerNumber{num, 0}
}

func (num IntegerNumber) IsNumber()        {}
func (num IntegerNumber) CanBeExpression() {}

func (num IntegerNumber) AddFloat(a float64) {
	num.Int += int64(a)
}

func (num IntegerNumber) AddInt(a int64) {
	num.Int += a
}

func (num IntegerNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num IntegerNumber) String() (out string) {
	out += fmt.Sprintf("%d", num.Int)

}
