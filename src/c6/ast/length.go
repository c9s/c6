package ast

import "strconv"

type Length struct {
	Value float64
	Unit  UnitType
	Token *Token
}

func (self Length) CanBeValue()      {}
func (self Length) CanBeExpression() {}
func (self Length) CanBeNode()       {}

func (self Length) String() (out string) {
	out += strconv.FormatFloat(self.Value, 'G', -1, 64)
	if self.Unit != UNIT_NONE {
		out += self.Unit.UnitString()
	}
	return out
}

func NewLength(val float64, unit UnitType, token *Token) *Length {
	return &Length{val, unit, token}
}

func LengthSubLength(a *Length, b *Length) *Length {
	var result = a.Value - b.Value
	return NewLength(result, UNIT_NONE, nil)
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func LengthDivLength(a *Length, b *Length) *Length {
	var result = a.Value / b.Value
	return NewLength(result, UNIT_NONE, nil)
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func LengthMulLength(a *Length, b *Length) *Length {
	var result = a.Value * b.Value
	return NewLength(result, UNIT_NONE, nil)
}

func LengthAddLength(a *Length, b *Length) *Length {
	var result = a.Value + b.Value
	return NewLength(result, UNIT_NONE, nil)
}
