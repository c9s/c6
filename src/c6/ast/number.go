package ast

import "strconv"

type Number struct {
	Value float64
	Token *Token
}

func NewNumberInt64(num int64, token *Token) *Number {
	return &Number{float64(num), token}
}

func NewNumber(num float64, token *Token) *Number {
	return &Number{num, token}
}

func (self Number) GetValueType() ValueType {
	return NumberValue
}

func (num Number) String() (out string) {
	return strconv.FormatFloat(num.Value, 'G', -1, 64)
}

func (num Number) Boolean() bool {
	return num.Value > 0
}
