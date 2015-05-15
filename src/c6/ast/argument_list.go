package ast

type Argument struct {
	VariableName *Token
	DefaultValue Expression
}

func NewArgumentWithToken(variableName *Token) *Argument {
	return &Argument{variableName, nil}
}
