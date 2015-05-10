package ast

type WhileStatement struct {
	Condition Expression
	Block     *Block
	ElseBlock *Block
}

func (stm WhileStatement) CanBeStatement() {}

func (stm WhileStatement) SetElseBlock(block *Block) {
	stm.ElseBlock = block
}

func (stm WhileStatement) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewWhileStatement(condition Expression, block *Block) *WhileStatement {
	return &WhileStatement{condition, block, nil}
}
