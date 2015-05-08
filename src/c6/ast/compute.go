package ast

type ComputableValue interface {
	GetValueType() ValueType
}

type ValueType uint16

const (
	NumberValue    ValueType = 0
	HexColorValue            = 1
	RGBAColorValue           = 2
	RGBColorValue            = 3
	ListValue                = 4
	MapValue                 = 5
)
