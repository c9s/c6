package ast

import "strconv"

type Length struct {
	Value float64
	Unit  *Unit
	Token *Token
}

func (self Length) CanBeNode() {}
func (self Length) GetValueType() ValueType {
	return LengthValue
}

func (self Length) Boolean() bool {
	return self.Value > 0
}

func (self Length) String() (out string) {
	out += strconv.FormatFloat(self.Value, 'G', -1, 64)
	if self.Unit != nil {
		out += self.Unit.String()
	}
	return out
}

func NewLength(val float64, unit *Unit, token *Token) *Length {
	return &Length{val, unit, token}
}
