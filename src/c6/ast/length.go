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
