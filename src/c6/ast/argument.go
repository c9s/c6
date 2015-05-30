package ast

import "strconv"

type Argument struct {
	Name           *Token
	DefaultValue   Expr
	VariableLength bool
	Position       int
}

func (arg Argument) String() string {
	// TODO: output DefaultValue, VariableLength, Position
	return arg.Name.Str + "(" + strconv.FormatInt(int64(arg.Position), 10) + ")"
}

func NewArgumentWithToken(name *Token) *Argument {
	return &Argument{
		Name:           name,
		DefaultValue:   nil,
		VariableLength: false,
		Position:       -1,
	}
}
