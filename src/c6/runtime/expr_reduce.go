package runtime

import "c6/ast"

func CanReduceExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		return !e.IsCssSlash()
	}
	return IsConstantExpression(expr)
}

/*
Reduce constant expression to constant.

@return (Value, ok)

ok = true means the expression is reduced to simple constant.
*/
func ReduceExpression(expr ast.Expression) (ast.Value, bool) {
	switch e := expr.(type) {
	case *ast.BinaryExpression:

		if exprLeft, ok := ReduceExpression(e.Left); ok {
			e.Left = exprLeft
		}
		if exprRight, ok := ReduceExpression(e.Right); ok {
			e.Right = exprRight
		}

	case *ast.UnaryExpression:

		if retExpr, ok := ReduceExpression(e.Expr); ok {
			e.Expr = retExpr
		}

	default:
		// it's already an constant value
		return e, true
	}

	if CanReduceExpression(expr) {
		switch e := expr.(type) {
		case *ast.BinaryExpression:
			return EvaluateBinaryExpression(e, nil), true
		case *ast.UnaryExpression:
			return EvaluateUnaryExpression(e, nil), true
		}
	}
	// not a constant expression
	return nil, false
}
