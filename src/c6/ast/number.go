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

/*
For color type, we treat them as vector. a vector can be computed with scalar or another vector.


Valid expression:

	#aaa / 3
	#bb * 3
	#bb - #cc
	3px - 3px
	3px + 3px

Invalid expression

- Number can't be the dividend.
- Hex Color can't divisor.

	3 / #aaa
	3 - #bbb
	3px - 3
	6em - 3px

*/

// Check the unit of the number operands to see if they're computable.
func NumberUnitCompatible(a *Number, b *Number) bool {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()
	// If any of them is without unit
	return unitA == UNIT_NONE || unitB == UNIT_NONE || unitA == unitB
}

func NumberSub(a *Number, b *Number) *Number {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()
	if unitA != unitB {
		panic("Can't compute number: incompatible number unit.")
	}

	var result = a.Value - b.Value
	var num = NewNumber(result, nil)
	num.SetUnit(unitA)
	return num
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func NumberDiv(a *Number, b *Number) *Number {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()

	// 10 / 3px is invalid
	if unitA == UNIT_NONE && unitB != UNIT_NONE {
		panic("Invalid number divisor")
	}

	var unitC UnitType = unitA

	// 10px / 10px = 1
	if unitA == unitB {
		unitC = UNIT_NONE
	}

	var result = a.Value / b.Value
	var num = NewNumber(result, nil)
	num.SetUnit(unitC)
	return num
	return nil
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func NumberMul(a *Number, b *Number) *Number {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()
	if !(unitA == UNIT_NONE || unitB == UNIT_NONE || unitA == unitB) {
		panic("Invalid number multiply expression")
	}

	var unitC UnitType = UNIT_NONE
	if unitA != UNIT_NONE {
		unitC = unitA
	} else if unitB != UNIT_NONE {
		unitC = unitB
	}

	var result = a.Value * b.Value
	var num = NewNumber(result, nil)
	num.SetUnit(unitC)
	return num
}

func NumberAdd(a *Number, b *Number) *Number {
	var unitA = a.GetUnit()
	var unitB = b.GetUnit()

	if unitA != unitB {
		panic("Can't compute number: incompatible number unit.")
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
