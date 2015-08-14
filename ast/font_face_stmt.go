package ast

type FontFaceStmt struct {
	Block *DeclBlock
}

func (stm FontFaceStmt) CanBeStmt()     {}
func (stm FontFaceStmt) String() string { return "" }
