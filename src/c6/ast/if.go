package ast

type If struct {
	Condition Expression
	Block     *Block
	ElseIf    []*If
}

func (stm If) CanBeStatement() {}

func (stm If) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIf() *If {
	return &If{}
}
