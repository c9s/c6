package ast

type WhileStmt struct {
	Condition Expr
	Block     *DeclarationBlock
	ElseBlock *DeclarationBlock
}

func (stm WhileStmt) CanBeStmt() {}

func (stm WhileStmt) SetElseBlock(block *DeclarationBlock) {
	stm.ElseBlock = block
}

func (stm WhileStmt) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewWhileStmt(condition Expr, block *DeclarationBlock) *WhileStmt {
	return &WhileStmt{condition, block, nil}
}
