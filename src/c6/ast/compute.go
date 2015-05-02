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
)

/*
Value
*/
type ComputeFunction func(a Value, b Value) Value

var computableMatrix [5][5]bool = [5][5]bool{
	/* NumberValue */
	[5]bool{false, false, false, false},

	/* LengthValue */
	[5]bool{false, false, false, false},

	/* HexColorValue */
	[5]bool{false, false, false, false},

	/* RGBAColorValue */
	[5]bool{false, false, false, false},

	/* RGBColorValue */
	[5]bool{false, false, false, false},
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
