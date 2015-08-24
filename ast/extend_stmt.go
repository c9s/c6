package ast

type ExtendStmt struct {
	Selectors *ComplexSelectorList
}

func (stm ExtendStmt) CanBeStmt()     {}
func (stm ExtendStmt) String() string { return "@extend" }

func NewExtendStmt() *ExtendStmt {
	return &ExtendStmt{}
}
