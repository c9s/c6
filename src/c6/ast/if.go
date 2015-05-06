package ast

type If struct {
	Condition Expression
	Block     *Block
	ElseIfs   []*If
}

func (stm If) CanBeStatement() {}

func (stm If) AppendElseIf(ifStm *If) {
	stm.ElseIfs = append(stm.ElseIfs, ifStm)
}

func (stm If) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIf() *If {
	return &If{}
}
