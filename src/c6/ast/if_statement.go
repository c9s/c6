package ast

type IfStatement struct {
	Condition Expression
	Block     *Block
	ElseIfs   []*IfStatement
	ElseBlock *Block
}

func (stm IfStatement) CanBeStatement() {}

func (stm IfStatement) AppendElseIf(ifStm *IfStatement) {
	stm.ElseIfs = append(stm.ElseIfs, ifStm)
}

func (stm IfStatement) SetElseBlock(block *Block) {
	stm.ElseBlock = block
}

func (stm IfStatement) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIfStatement(condition Expression, block *Block) *IfStatement {
	return &IfStatement{condition, block, []*IfStatement{}, nil}
}
