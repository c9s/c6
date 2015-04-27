package ast

type List struct {
	Separator   string
	Expressions []Expression
}

func (list List) CanBeExpression() {}
func (list List) CanBeValue()      {}

func (list *List) Append(expr Expression) {
	list.Expressions = append(list.Expressions, expr)
}

// By the default, the separator is space
func NewList() *List {
	return &List{" ", []Expression{}}
}
