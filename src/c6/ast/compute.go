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

var computableValueMatrix1 [5][5]bool = [5][5]bool{
	/* NumberValue */
	[5]bool{false, false, false, false},

	/* HexColorValue */
	[5]bool{false, false, false, false},

	/* RGBAColorValue */
	[5]bool{false, false, false, false},

	/* RGBColorValue */
	[5]bool{false, false, false, false},
}

var computableValueMatrix map[ValueType]map[ValueType]bool = map[ValueType]map[ValueType]bool{
	NumberValue: map[ValueType]bool{
		NumberValue:    true,
		HexColorValue:  true,
		RGBAColorValue: true,
		RGBColorValue:  true,
	},
	HexColorValue: map[ValueType]bool{
		NumberValue:    true,
		HexColorValue:  true,
		RGBAColorValue: false,
		RGBColorValue:  false,
	},
	RGBAColorValue: map[ValueType]bool{
		NumberValue:    true,
		HexColorValue:  false,
		RGBAColorValue: true,
		RGBColorValue:  false,
	},
	RGBColorValue: map[ValueType]bool{
		NumberValue:    true,
		HexColorValue:  false,
		RGBAColorValue: false,
		RGBColorValue:  true,
	},
}
