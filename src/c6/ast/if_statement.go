package ast

type IfStatement struct {
	Condition Expression
	Block     *DeclarationBlock
	ElseIfs   []*IfStatement
	ElseBlock *DeclarationBlock
}

func (stm IfStatement) CanBeStatement() {}

func (stm IfStatement) AppendElseIf(ifStm *IfStatement) {
	stm.ElseIfs = append(stm.ElseIfs, ifStm)
}

func (stm IfStatement) SetElseBlock(block *DeclarationBlock) {
	stm.ElseBlock = block
}

func (stm IfStatement) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIfStatement(condition Expression, block *DeclarationBlock) *IfStatement {
	return &IfStatement{condition, block, []*IfStatement{}, nil}
}
