package ast

import "strings"

type List struct {
	Separator   string
	Expressions []Expression
}

func (self List) GetValueType() ValueType {
	return ListValue
}

func (list List) String() string {
	var exprstrs []string
	for _, expr := range list.Expressions {
		exprstrs = append(exprstrs, expr.String())
	}
	return strings.Join(exprstrs, list.Separator)
}

func (list *List) Len() int {
	return len(list.Expressions)
}

func (list *List) Append(expr Expression) {
	var newList = append(list.Expressions, expr)
	list.Expressions = newList
}

// By the default, the separator is space
func NewList(sep string) *List {
	return &List{sep, []Expression{}}
}

func NewSpaceSepList() *List {
	return &List{" ", []Expression{}}
}

func NewCommaSepList() *List {
	return &List{", ", []Expression{}}
}
