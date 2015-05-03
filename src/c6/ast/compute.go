package ast

type ComputableValue interface {
	GetValueType() ValueType
}

type ValueType uint16

const (
	NumberValue    ValueType = 0
	LengthValue              = 1
	HexColorValue            = 2
	RGBAColorValue           = 3
	RGBColorValue            = 4
	ListValue                = 5
	MapValue                 = 6
)

/*
Value
*/
type ComputeFunction func(a Value, b Value) Value

const ValueTypeNum = 7

var computableMatrix [ValueTypeNum][ValueTypeNum]bool = [ValueTypeNum][ValueTypeNum]bool{
	/* NumberValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* LengthValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* HexColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBAColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},
}

/**
Each row: [5]ComputeFunction{ NumberValue, LengthValue, HexColorValue, RGBAColorValue, RGBColorValue }
*/
var computeFunctionMatrix [5][5]ComputeFunction = [5][5]ComputeFunction{

	/* NumberValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* LengthValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* HexColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBAColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},
}

func Compute(op OpType, a Value, b Value) Value {
	switch op {
	case OpAdd:
		switch ta := a.(type) {
		case *Number:
			switch tb := b.(type) {
			case *Number:
				return NumberAddNumber(ta, tb)
			case *HexColor:
				return HexColorAddNumber(tb, ta)
			}
		case *Length:
			switch tb := b.(type) {
			case *Length:
				return LengthAddLength(ta, tb)
			}
		case *HexColor:
			switch tb := b.(type) {
			case *Number:
				return HexColorAddNumber(ta, tb)
			}
		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorAddNumber(ta, tb)
			}
		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorAddNumber(ta, tb)
			}
		}
	case OpSub:
		switch ta := a.(type) {

		case *Number:
			switch tb := b.(type) {
			case *Number:
				return NumberSubNumber(ta, tb)
			}

		case *Length:
			switch tb := b.(type) {
			case *Length:
				val := LengthSubLength(ta, tb)
				fmt.Printf("Substracted value: %+v\n", val)
				return val
			}

		case *HexColor:
			switch tb := b.(type) {
			case *Number:
				return HexColorSubNumber(ta, tb)
			}

		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorSubNumber(ta, tb)
			}

		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorSubNumber(ta, tb)
			}
		}
	case OpMul:
		switch ta := a.(type) {

		case *Length:
			switch tb := b.(type) {
			case *Length:
				return LengthMulLength(ta, tb)
			case *Number:
				return LengthMulNumber(ta, tb)
			}

		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorMulNumber(ta, tb)
			}

		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorMulNumber(ta, tb)
			}
		}
	}
	return nil
}
