package ast

type Interpolation struct {
	Expression Expression
}

func (self Interpolation) CanBeExpression() {}
func (self Interpolation) CanBeNode()       {}
