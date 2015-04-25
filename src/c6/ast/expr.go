package ast

type Expression interface {
	CanBeExpression()
}

type UnaryExpression struct {
	Op    Op
	Expr  Expression
	Token Token
}

func NewUnaryExpression(op Op, expr Expression, token Token) *UnaryExpression {
	return &UnaryExpression{op, expr, token}
}

func (self UnaryExpression) CanBeExpression() {}

type BinaryExpression struct {
	Op    Op
	Left  Expression
	Right Expression
}

func (self BinaryExpression) CanBeExpression() {}

func NewBinaryExpression(op Op, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}
