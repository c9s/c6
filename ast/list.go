package ast

import "strings"

type List struct {
	Separator string
	Exprs     []Expr
}

/*
When there is more than one item, we return true for boolean context
*/
func (self List) Boolean() bool {
	return len(self.Exprs) > 0
}

func (self List) GetValueType() ValueType {
	return ListValue
}

func (list List) String() string {
	var exprstrs []string
	for _, expr := range list.Exprs {
		exprstrs = append(exprstrs, expr.String())
	}
	return strings.Join(exprstrs, list.Separator)
}

func (list *List) Len() int {
	return len(list.Exprs)
}

func (list *List) Append(expr Expr) {
	var newList = append(list.Exprs, expr)
	list.Exprs = newList
}

// By the default, the separator is space
func NewList(sep string) *List {
	return &List{sep, []Expr{}}
}

func NewSpaceSepList() *List {
	return &List{" ", []Expr{}}
}

func NewCommaSepList() *List {
	return &List{", ", []Expr{}}
}
