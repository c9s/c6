package ast

type Expression interface {
	CanBeExpression()
	String() string
}

type UnaryExpression struct {
	Op   *Op
	Expr Expression
}

func NewUnaryExpression(op *Op, expr Expression) *UnaryExpression {
	return &UnaryExpression{op, expr}
}

func (self UnaryExpression) CanBeExpression() {}
func (self UnaryExpression) String() string {
	if self.Op != nil {
		return self.Op.String() + self.Expr.String()
	}
	return self.Expr.String()
}

type BinaryExpression struct {
	Op    *Op
	Left  Expression
	Right Expression
}

func (self BinaryExpression) CanBeExpression() {}

func (self BinaryExpression) String() string {
	return self.Left.String() + self.Op.String() + self.Right.String()
}

func NewBinaryExpression(op *Op, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}
