package ast

type Number interface {
	SetUnit(unit int)
	IsNumber()
	CanBeExpression()
}

type FloatNumber struct {
	Float float64
	Unit  int
}

func NewFloatNumber(num float64) *FloatNumber {
	return &FloatNumber{num, 0}
}

func (num FloatNumber) IsNumber()        {}
func (num FloatNumber) CanBeExpression() {}

func (num FloatNumber) SetUnit(unit int) {
	num.Unit = unit
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
func (num IntegerNumber) SetUnit(unit int) {
	num.Unit = unit
}
