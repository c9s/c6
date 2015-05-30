package runtime

import "c6/ast"

func CanReduceExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return !e.IsCssSlash()
	}
	return true
}

/*
Reduce constant expression to constant.

@return (Value, ok)

ok = true means the expression is reduced to simple constant.

The difference between Evaluate*Expr method is:

- `ReduceExpr` returns either value or expression (when there is an unsolved expression)
- `EvaluateBinaryExpr` returns nil if there is an unsolved expression.

*/
func ReduceExpr(expr ast.Expr, context *Context) (ast.Value, bool) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if exprLeft, ok := ReduceExpr(e.Left, context); ok {
			e.Left = exprLeft
		}
		if exprRight, ok := ReduceExpr(e.Right, context); ok {
			e.Right = exprRight
		}

	case *ast.UnaryExpr:

		if retExpr, ok := ReduceExpr(e.Expr, context); ok {
			e.Expr = retExpr
		}

	case *ast.Variable:
		if context == nil {
			return nil, false
		}

		if varVal, ok := context.GetVariable(e.Name); ok {
			return varVal.(ast.Expr), true
		}

	default:
		// it's already an constant value
		return e, true
	}

	if IsSimpleExpr(expr) {
		switch e := expr.(type) {
		case *ast.BinaryExpr:
			return EvaluateBinaryExpr(e, context), true
		case *ast.UnaryExpr:
			return EvaluateUnaryExpr(e, context), true
		}
	}
	// not a constant expression
	return nil, false
}
