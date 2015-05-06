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
