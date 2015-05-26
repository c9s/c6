package c6

import "c6/ast"

func CanReduceExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		return !e.IsCssSlash()
	}
	return true
}

/*
Reduce constant expression to constant.

@return (Value, ok)

ok = true means the expression is reduced to simple constant.

The difference between Evaluate*Expression method is:

- `ReduceExpression` returns either value or expression (when there is an unsolved expression)
- `EvaluateBinaryExpression` returns nil if there is an unsolved expression.

*/
func ReduceExpression(expr ast.Expression, context *Context) (ast.Value, bool) {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		if exprLeft, ok := ReduceExpression(e.Left, context); ok {
			e.Left = exprLeft
		}
		if exprRight, ok := ReduceExpression(e.Right, context); ok {
			e.Right = exprRight
		}

	case *ast.UnaryExpression:

		if retExpr, ok := ReduceExpression(e.Expr, context); ok {
			e.Expr = retExpr
		}

	case *ast.Variable:
		if context == nil {
			return nil, false
		}

		if varVal, ok := context.GetVariable(e.Name); ok {
			return varVal.(ast.Expression), true
		}

	default:
		// it's already an constant value
		return e, true
	}

	if IsSimpleExpression(expr) {
		switch e := expr.(type) {
		case *ast.BinaryExpression:
			return EvaluateBinaryExpression(e, context), true
		case *ast.UnaryExpression:
			return EvaluateUnaryExpression(e, context), true
		}
	}
	// not a constant expression
	return nil, false
}
