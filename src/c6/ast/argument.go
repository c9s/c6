package ast

type Argument struct {
	VariableName   *Token
	DefaultValue   Expression
	VariableLength bool
}

func NewArgumentWithToken(variableName *Token) *Argument {
	return &Argument{variableName, nil, false}
}

type ArgumentList []*Argument

func (args ArgumentList) Append(arg *Argument) {
	newargs := append(args, arg)
	args = newargs
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{}
}
