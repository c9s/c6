package ast

/**
A struct that contains many concatable strings

	{any} {expression} {any}
*/
type Concat struct {
	Expressions []Expression
}

func NewConcat() *Concat {
	return &Concat{[]Expression{}}
}

func (self Concat) AppendExpression(expr Expression) {
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
