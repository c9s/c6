package ast

type IfStatement struct {
	Condition Expression
	Block     *Block
	ElseIfs   []*IfStatement
}

func (stm IfStatement) CanBeStatement() {}

func (stm IfStatement) AppendElseIf(ifStm *IfStatement) {
	stm.ElseIfs = append(stm.ElseIfs, ifStm)
}

func (stm IfStatement) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIfStatement() *IfStatement {
	return &IfStatement{}
}
