package ast

import "fmt"

type Number interface {
	SetUnit(unit int)
	IsNumber()
	CanBeExpression()
}

type FloatNumber struct {
	Float float64
	Unit  int
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

func (num FloatNumber) SetUnit(unit int) {
	num.Unit = unit
}

func (num FloatNumber) String() (out string) {
	out += fmt.Sprintf("%.2f", num.Float)

	if num.Unit > 0 {
		switch TokenType(num.Unit) {
		case T_UNIT_PX:
			out += "px"
		case T_UNIT_PT:
			out += "pt"
		case T_UNIT_EM:
			out += "em"
		case T_UNIT_REM:
			out += "rem"
		default:
			panic("unimplemented type")
		}
	}
	return out
}

type IntegerNumber struct {
	Int  int64
	Unit int
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

func (num IntegerNumber) SetUnit(unit int) {
	num.Unit = unit
}
