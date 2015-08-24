package ast

type WhileStmt struct {
	Condition Expr
	Block     *DeclBlock
	ElseBlock *DeclBlock
}

func (stm WhileStmt) CanBeStmt() {}

func (stm WhileStmt) SetElseBlock(block *DeclBlock) {
	stm.ElseBlock = block
}

func (stm WhileStmt) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewWhileStmt(condition Expr, block *DeclBlock) *WhileStmt {
	return &WhileStmt{condition, block, nil}
}
