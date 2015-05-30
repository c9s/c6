package runtime

import "c6/ast"

type Stack struct {
	Values []ast.Expr
}

func (stack *Stack) Push(arg ast.Expr) {
	stack.Values = append(stack.Values, arg)
}

func (stack *Stack) Pop() ast.Expr {
	if len(stack.Values) > 0 {
		var val = stack.Values[len(stack.Values)-1]
		stack.Values = stack.Values[:len(stack.Values)-1]
		return val
	}
	return nil
}
