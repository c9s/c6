package ast

type Expression interface {
	CanBeExpression()
}

type UnaryExpression struct {
	Value interface{}
	Token Token
}

func (self UnaryExpression) CanBeExpression() {}

type BinaryExpression struct {
	Left  Expression
	Right Expression
	Op    string
}

func (self BinaryExpression) CanBeExpression() {}

type ConstantString struct {
	Constant string
	Token    Token
}
