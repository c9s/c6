package ast

/**
A struct that contains many concatable strings

	{any}{expression}{any}
*/
type LiteralConcat struct {
	Expressions []Expression
}

func NewLiteralConcat() *LiteralConcat {
	return &LiteralConcat{[]Expression{}}
}

func (self LiteralConcat) String() (out string) {
	out = ""
	for _, expr := range self.Expressions {
		out += expr.String()
	}
	return out
}

func (self LiteralConcat) CanBeExpression() {}

func (self *LiteralConcat) AppendExpression(expr Expression) {
	self.Expressions = append(self.Expressions, expr)
}

/*
func (self Concat) String() string {
	var out string = ""

	for _, expr := range self.Expressions {
		_ = expr
		// expr.String()
	}
	return out
}
*/
