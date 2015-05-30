package ast

/**
A struct that contains many concatable strings

	{any}{expression}{any}
*/
type LiteralConcat struct {
	Left  Expr
	Right Expr
}

func NewLiteralConcat(left, right Expr) *LiteralConcat {
	return &LiteralConcat{left, right}
}

func (self LiteralConcat) String() string {
	return self.Left.String() + self.Right.String()
}
