package ast

type IfStmt struct {
	Condition Expr
	Block     *DeclarationBlock
	ElseIfs   []*IfStmt
	ElseBlock *DeclarationBlock
}

func (stm IfStmt) CanBeStmt() {}

func (stm IfStmt) AppendElseIf(ifStm *IfStmt) {
	stm.ElseIfs = append(stm.ElseIfs, ifStm)
}

func (stm IfStmt) SetElseBlock(block *DeclarationBlock) {
	stm.ElseBlock = block
}

func (stm IfStmt) String() string {
	return "(if statement STRING() un-implemented)"
}

func NewIfStmt(condition Expr, block *DeclarationBlock) *IfStmt {
	return &IfStmt{condition, block, []*IfStmt{}, nil}
}
