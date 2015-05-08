package ast

import "strconv"

type Number struct {
	Value  float64
	double bool
	Unit   *Unit
	Token  *Token
}

func NewNumber(num float64, unit *Unit, token *Token) *Number {
	return &Number{num, false, unit, token}
}

/*
Mark the number as an double (value with precision)
*/
func (num *Number) SetDouble() {
	num.double = true
}

/*
Check if the number is a floating value.
*/
func (num *Number) IsDouble() bool {
	return num.double
}

func (self Number) GetValueType() ValueType {
	return NumberValue
}

func (self Number) String() (out string) {
	out += strconv.FormatFloat(self.Value, 'G', -1, 64)
	if self.Unit != nil {
		out += self.Unit.String()
	}
	return out
}

func (num Number) Double() float64 {
	return num.Value
}

func (num Number) Integer() int {
	return int(num.Value)
}

func (num Number) Boolean() bool {
	return num.Value > 0
}
