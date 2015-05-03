package ast

import "strconv"

type Length struct {
	Value float64
	Unit  UnitType
	Token *Token
}

func (self Length) CanBeNode() {}
func (self Length) GetValueType() ValueType {
	return LengthValue
}

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
	if a.Unit != b.Unit {
		panic("Incompatable unit")
	}
	var result = a.Value - b.Value
	return NewLength(result, UNIT_NONE, nil)
}

func LengthAddLength(a *Length, b *Length) *Length {
	if a.Unit != b.Unit {
		panic("Incompatable unit")
	}
	var result = a.Value + b.Value
	return NewLength(result, UNIT_NONE, nil)
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func LengthDivLength(a *Length, b *Length) *Length {
	if a.Unit == UNIT_NONE || b.Unit == UNIT_NONE || a.Unit == b.Unit {
		var result = a.Value / b.Value
		var unit = UNIT_NONE
		if a.Unit != UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit != UNIT_NONE {
			unit = b.Unit
		}
		return NewLength(result, unit, nil)
	}
	panic("Incompatable unit")
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func LengthMulLength(a *Length, b *Length) *Length {
	if a.Unit == UNIT_NONE || b.Unit == UNIT_NONE || a.Unit == b.Unit {
		var result = a.Value * b.Value
		var unit = UNIT_NONE
		if a.Unit != UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit != UNIT_NONE {
			unit = b.Unit
		}
		return NewLength(result, unit, nil)
	}
	panic("Incompatable unit")
}

func LengthMulNumber(a *Length, b *Number) *Length {
	return NewLength(a.Value*b.Value, a.Unit, nil)
}

func NumberMulLength(a *Number, b *Length) *Length {
	return NewLength(a.Value*b.Value, b.Unit, nil)
}

func LengthDivNumber(a *Length, b *Number) *Length {
	return NewLength(a.Value/b.Value, a.Unit, nil)
}

func NumberDivLength(a *Number, b *Length) *Length {
	panic("Number can't be divided by length")
	return nil
}
