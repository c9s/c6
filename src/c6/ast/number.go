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

func (num Number) CanBeValue()      {}
func (num Number) CanBeNode()       {}
func (num Number) CanBeExpression() {}

func (self Number) GetValueType() ValueType {
	return NumberValue
}

func (num Number) String() (out string) {
	return strconv.FormatFloat(num.Value, 'G', -1, 64)
}

/*
For color type, we treat them as vector. a vector can be computed with scalar or another vector.


Valid expression:

	#aaa / 3
	#bb * 3
	#bb - #cc
	3px - 3px
	3px + 3px

Invalid expression

- Number can't be the dividend.
- Hex Color can't divisor.

	3 / #aaa
	3 - #bbb
	3px - 3
	6em - 3px

*/

func NumberSubNumber(a *Number, b *Number) *Number {
	var result = a.Value - b.Value
	var num = NewNumber(result, nil)
	return num
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func NumberDivNumber(a *Number, b *Number) *Number {
	var result = a.Value / b.Value
	return NewNumber(result, nil)
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func NumberMulNumber(a *Number, b *Number) *Number {
	var result = a.Value * b.Value
	return NewNumber(result, nil)
}

func NumberAddNumber(a *Number, b *Number) *Number {
	var result = a.Value + b.Value
	return NewNumber(result, nil)
}
