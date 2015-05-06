package ast

import "strconv"

type BooleanValue interface {
	Boolean() bool
}

type Boolean struct {
	Value bool
	Token *Token
}

func (self Boolean) Boolean() bool {
	return self.Value
}

func (self Boolean) String() string {
	if self.Token != nil {
		return self.Token.Str
	}
	return strconv.FormatBool(self.Value)
}

func NewBooleanTrue(token *Token) *Boolean {
	return &Boolean{true, token}
}

func NewBooleanFalse(token *Token) *Boolean {
	return &Boolean{false, token}
}

func NewBoolean(val bool) *Boolean {
	return &Boolean{val, nil}
}

func NewBooleanWithToken(token *Token) *Boolean {
	val, err := strconv.ParseBool(token.Str)
	if err != nil {
		panic("Can't parse boolean value")
	}
	return &Boolean{val, token}
}
