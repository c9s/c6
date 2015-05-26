package ast

type Argument struct {
	Name           *Token
	DefaultValue   Expression
	VariableLength bool
	Position       uint8
}

func NewArgumentWithToken(name *Token) *Argument {
	return &Argument{
		Name:           name,
		DefaultValue:   nil,
		VariableLength: false,
	}
}
