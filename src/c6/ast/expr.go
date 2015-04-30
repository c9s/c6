package ast

type Expression interface {
	CanBeExpression()
	String() string
}

type UnaryExpression struct {
	Op   OpType
	Expr Expression
}

func NewUnaryExpression(op OpType, expr Expression) *UnaryExpression {
	return &UnaryExpression{op, expr}
}

func (self UnaryExpression) CanBeExpression() {}
func (self UnaryExpression) String() string {
	return "(" + self.Op.String() + " " + self.Expr.String() + ")"
}

type BinaryExpression struct {
	Op    OpType
	Left  Expression
	Right Expression
}

func (self BinaryExpression) CanBeExpression() {}

func (self BinaryExpression) String() string {
	return "(" + self.Left.String() + " " + self.Op.String() + " " + self.Right.String() + ")"
}

func NewBinaryExpression(op OpType, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}
