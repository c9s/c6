package ast

type FontFaceStatement struct {
	Block *DeclarationBlock
}

func (stm FontFaceStatement) CanBeStatement() {}
func (stm FontFaceStatement) String() string  { return "" }
