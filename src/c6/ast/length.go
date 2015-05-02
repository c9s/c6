package ast

type Length struct {
	Value float64
	Unit  UnitType
	Token *Token
}

func NewLength(val float64, unit UnitType, token *Token) *Length {
	return &Length{val, unit, token}
}
