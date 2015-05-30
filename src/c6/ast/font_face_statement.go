package ast

type FontFaceStmt struct {
	Block *DeclarationBlock
}

func (stm FontFaceStmt) CanBeStmt() {}
func (stm FontFaceStmt) String() string  { return "" }
