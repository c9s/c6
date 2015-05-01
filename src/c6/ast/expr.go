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

func (self BinaryExpression) String() (out string) {
	out += "("
	out += self.Left.String()
	out += " " + self.Op.String() + " "
	out += self.Right.String()
	out += ")"
	return out
}

func NewBinaryExpression(op OpType, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}
