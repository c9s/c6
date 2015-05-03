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
