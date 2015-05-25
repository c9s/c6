package ast

type ArgumentList []*Argument

func (args ArgumentList) Append(arg *Argument) {
	newargs := append(args, arg)
	args = newargs
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{}
}
