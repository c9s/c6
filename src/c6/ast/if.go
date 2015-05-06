package ast

type If struct {
	Condition Node
	Block
	ElseIf []If
}

func NewIf() *If {
	return &If{}
}
