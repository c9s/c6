package ast

type WhileStatement struct {
	Condition Expression
	Block     *DeclarationBlock
	ElseBlock *DeclarationBlock
}

func (stm WhileStatement) CanBeStatement() {}

func (stm WhileStatement) SetElseBlock(block *DeclarationBlock) {
	stm.ElseBlock = block
}

func (stm WhileStatement) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewWhileStatement(condition Expression, block *DeclarationBlock) *WhileStatement {
	return &WhileStatement{condition, block, nil}
}
