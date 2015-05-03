package ast

type Expression interface {
	CanBeExpression()
	// Evaluate() Value
	String() string
}

type UnaryExpression struct {
	Op   OpType
	Expr Expression
}

func NewUnaryExpression(op OpType, expr Expression) *UnaryExpression {
	return &UnaryExpression{op, expr}
}

func (self *UnaryExpression) Evaluate() Value {
	if self.Op == OpAdd {
		// TODO: call Evaluate()
		return Value(self.Expr)
	} else if self.Op == OpSub {
		// TODO: call Evaluate()
		return Value(self.Expr)
	}
	return nil
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

func (self *BinaryExpression) Evaluate() Value {
	if self.Op == OpAdd {
		// TODO: call Evaluate()
		return Value(self.Left)
	} else if self.Op == OpSub {
		// TODO: call Evaluate()
		return Value(self.Left)
	}
	return nil
}

func NewBinaryExpression(op OpType, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}
