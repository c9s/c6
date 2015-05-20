package ast

type ExtendStatement struct {
	Selectors *ComplexSelectorList
}

func (stm ExtendStatement) CanBeStatement() {}
func (stm ExtendStatement) String() string  { return "@extend" }

func NewExtendStatement() *ExtendStatement {
	return &ExtendStatement{}
}
