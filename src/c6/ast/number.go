package ast

type Number interface {
	SetUnit(unit UnitType)
	GetUnit() UnitType
	CanBeNumber()
	CanBeExpression()
}

func NewNumber(val interface{}) Number {
	floatval, ok := val.(float64)
	if ok {
		return NewFloatNumber(floatval)
	}
	intval, ok := val.(int64)
	if ok {
		return NewIntegerNumber(intval)
	}
	panic("Unknown number type")
}

// Pass numbers as pointer
func AddNumber(a Number, b Number) Number {
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

	switch a.(type) {
	case *FloatNumber:
		var result float64 = a.(*FloatNumber).Float
		switch b.(type) {
		case *IntegerNumber:
			result += float64(b.(*IntegerNumber).Int)
		case *FloatNumber:
			result += float64(b.(*FloatNumber).Float)
		default:
			panic("Unknown type for calculation")
		}
		num := NewNumber(result)
		num.SetUnit(unitC)
		return num
	case *IntegerNumber:
		switch b.(type) {
		case *IntegerNumber:
			var result int64 = a.(*IntegerNumber).Int
			result += int64(b.(*IntegerNumber).Int)
			var num = NewNumber(result)
			num.SetUnit(unitC)
			return num
		case *FloatNumber:
			var result float64 = float64(a.(*IntegerNumber).Int)
			result += float64(b.(*FloatNumber).Float)
			var num = NewNumber(result)
			num.SetUnit(unitC)
			return num
		default:
			panic("Unknown type for calculation")
		}
	}
	panic("Unknown type for calculation")
	return nil
}
