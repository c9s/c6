package ast

type Argument struct {
	VariableName   *Token
	DefaultValue   Expression
	VariableLength bool
}

func NewArgumentWithToken(variableName *Token) *Argument {
	return &Argument{variableName, nil, false}
}
