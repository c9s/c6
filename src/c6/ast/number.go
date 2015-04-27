package ast

import "fmt"
import "strconv"

type Number interface {
	SetUnit(unit UnitType)
	GetUnit() UnitType
	CanBeNumber()
	CanBeExpression()
}

type FloatNumber struct {
	Float float64
	Unit  UnitType
	Token *Token
}

func (num *FloatNumber) AddFloat(a float64) {
	num.Float += a
}

func (num *FloatNumber) AddInt(a int) {
	num.Float += float64(a)
}

func (num *FloatNumber) GetUnit() UnitType {
	return num.Unit
}

func NewFloatNumber(num float64) *FloatNumber {
	return &FloatNumber{num, UNIT_NONE, nil}
}

func (num FloatNumber) CanBeNumber()     {}
func (num FloatNumber) CanBeNode()       {}
func (num FloatNumber) CanBeExpression() {}

func (num *FloatNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num FloatNumber) String() (out string) {
	out += strconv.FormatFloat(num.Float, 'G', -1, 64)
	if num.Unit != UNIT_NONE {
		out += num.Unit.UnitString()
	}
	return out
}

type IntegerNumber struct {
	Int   int64
	Unit  UnitType
	Token *Token
}

func NewIntegerNumber(num int64) *IntegerNumber {
	return &IntegerNumber{num, UNIT_NONE, nil}
}

func (num IntegerNumber) CanBeNumber()     {}
func (num IntegerNumber) CanBeNode()       {}
func (num IntegerNumber) CanBeExpression() {}

func (num *IntegerNumber) GetUnit() UnitType {
	return num.Unit
}

func (num *IntegerNumber) AddFloat(a float64) {
	num.Int += int64(a)
}

func (num *IntegerNumber) AddInt(a int64) {
	num.Int += a
}

func (num *IntegerNumber) SetUnit(unit UnitType) {
	num.Unit = unit
}

func (num IntegerNumber) String() (out string) {
	out += fmt.Sprintf("%d", num.Int)
	if num.Unit > 0 {
		out += num.Unit.UnitString()
	}
	return out
}

type NumberType int

const (
	Float NumberType = iota
	Integer
)

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
