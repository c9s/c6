package ast

import "strconv"

type Number struct {
	Value float64
	Unit  UnitType
	Token *Token
}

func (num *Number) AddFloat(a float64) {
	num.Value += a
}

func (num *Number) AddInt(a int) {
	num.Value += float64(a)
}

func (num *Number) GetUnit() UnitType {
	return num.Unit
}

func NewNumberInt64(num int64, token *Token) *Number {
	return &Number{float64(num), UNIT_NONE, token}
}

func NewNumber(num float64, token *Token) *Number {
	return &Number{num, UNIT_NONE, token}
}

func (num Number) CanBeValue()      {}
func (num Number) CanBeNode()       {}
func (num Number) CanBeExpression() {}

func (num *Number) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num Number) String() (out string) {
	out += strconv.FormatFloat(num.Value, 'G', -1, 64)
	if num.Unit != UNIT_NONE {
		out += num.Unit.UnitString()
	}
	return out
}

// Pass numbers as pointer
func AddNumber(a *Number, b *Number) *Number {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()
	if unitA != UNIT_NONE && unitB != UNIT_NONE && a.GetUnit() != b.GetUnit() {
		panic("Incompatible number type")
	}

	var unitC UnitType = UNIT_NONE
	if unitA != UNIT_NONE {
		unitC = unitA
	} else if unitB != UNIT_NONE {
		unitC = unitB
	}

	var result = a.Value + b.Value
	var num = NewNumber(result, nil)
	num.SetUnit(unitC)
	return num
}
