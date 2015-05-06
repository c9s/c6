package ast

import "strconv"

type Number struct {
	Value float64
	Unit  *Unit
	Token *Token
}

func NewNumber(num float64, unit *Unit, token *Token) *Number {
	return &Number{num, unit, token}
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

func (num Number) Boolean() bool {
	return num.Value > 0
}
