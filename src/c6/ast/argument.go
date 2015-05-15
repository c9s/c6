package ast

type Argument struct {
	VariableName *Token
	DefaultValue Expression
}

func NewArgumentWithToken(variableName *Token) *Argument {
	return &Argument{variableName, nil}
}

type ArgumentList []*Argument

func (args ArgumentList) Append(arg *Argument) {
	newargs := append(args, arg)
	args = newargs
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{}
}
