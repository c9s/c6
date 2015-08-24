package runtime

import "github.com/c9s/c6/ast"
import "fmt"

/*
Used for Incompatible unit, data type or unsupported operations

TODO: This is not used yet. our compute functions should return error if possible
*/
type ComputeError struct {
	Message string
	Left    ast.Value
	Right   ast.Value
}

func (self ComputeError) Error() string {
	return self.Message
}

/*
Value
*/
type ComputeFunction func(a ast.Value, b ast.Value) ast.Value

const ValueTypeNum = 7

var computableMatrix [ValueTypeNum][ValueTypeNum]bool = [ValueTypeNum][ValueTypeNum]bool{
	/* NumberValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* HexColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBAColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},
}

/**
Each row: [5]ComputeFunction{ NumberValue, HexColorValue, RGBAColorValue, RGBColorValue }
*/
var computeFunctionMatrix [5][5]ComputeFunction = [5][5]ComputeFunction{

	/* NumberValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* HexColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBAColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},
}

func Compute(op *ast.Op, a ast.Value, b ast.Value) ast.Value {
	if op == nil {
		panic("op can't be nil")
	}
	switch op.Type {

	case ast.T_EQUAL:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value == tb.Value)
			}
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value == tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_UNEQUAL:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value != tb.Value)
			}
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value != tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_GT:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value > tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_GE:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value >= tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LT:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value < tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LE:

		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				if IsComparable(ta, tb) {
					return ast.NewBoolean(ta.Value <= tb.Value)
				} else {
					panic("Can't compare number (unit different)")
				}
			}
		}

	case ast.T_LOGICAL_AND:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {

			case *ast.Boolean:
				return ast.NewBoolean(ta.Value && tb.Value)

			// For other data type, we cast to boolean
			default:
				if bv, ok := b.(ast.BooleanValue); ok {
					return ast.NewBoolean(bv.Boolean())
				}
			}
		}

	case ast.T_LOGICAL_OR:

		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {

			case *ast.Boolean:
				return ast.NewBoolean(ta.Value || tb.Value)

			// For other data type, we cast to boolean
			default:
				if bv, ok := b.(ast.BooleanValue); ok {
					return ast.NewBoolean(bv.Boolean())
				}
			}
		}

	/*
		arith expr
	*/
	case ast.T_PLUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberAddNumber(ta, tb)
			case *ast.HexColor:
				return HexColorAddNumber(tb, ta)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorAddNumber(ta, tb)
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorAddNumber(ta, tb)
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorAddNumber(ta, tb)
			}
		}
	case ast.T_MINUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberSubNumber(ta, tb)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorSubNumber(ta, tb)
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorSubNumber(ta, tb)
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorSubNumber(ta, tb)
			}
		}

	case ast.T_DIV:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberDivNumber(ta, tb)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorDivNumber(ta, tb)
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorDivNumber(ta, tb)
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorDivNumber(ta, tb)
			}
		}

	case ast.T_MUL:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberMulNumber(ta, tb)
			}

		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorMulNumber(ta, tb)
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorMulNumber(ta, tb)
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorMulNumber(ta, tb)
			}
		}
	}
	return nil
}

/*
A simple expression means the operands are scalar, and can be evaluated.
*/
func IsSimpleExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if IsValue(e.Left) && IsValue(e.Right) {
			return true
		}
	case *ast.UnaryExpr:
		if IsValue(e.Expr) {
			return true
		}
	}
	return false
}

/*
This function returns true when the val is a scalar value, not an expression.
*/
func IsValue(val ast.Expr) bool {
	switch val.(type) {
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor, *ast.HSLColor, *ast.HSVColor, *ast.Boolean:
		return true
	}
	return false
}

func EvaluateExprInBooleanContext(anyexpr ast.Expr, context *Context) ast.Value {
	switch expr := anyexpr.(type) {

	case *ast.BinaryExpr:
		return EvaluateBinaryExprInBooleanContext(expr, context)

	case *ast.UnaryExpr:
		return EvaluateUnaryExprInBooleanContext(expr, context)

	default:
		if bval, ok := expr.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean())
		}
	}
	return nil
}

func EvaluateBinaryExprInBooleanContext(expr *ast.BinaryExpr, context *Context) ast.Value {

	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {
	case *ast.UnaryExpr:
		lval = EvaluateUnaryExprInBooleanContext(expr, context)

	case *ast.BinaryExpr:
		lval = EvaluateBinaryExprInBooleanContext(expr, context)

	default:
		lval = expr
	}

	switch expr := expr.Right.(type) {
	case *ast.UnaryExpr:
		rval = EvaluateUnaryExprInBooleanContext(expr, context)

	case *ast.BinaryExpr:
		rval = EvaluateBinaryExprInBooleanContext(expr, context)

	default:
		rval = expr
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExprInBooleanContext(expr *ast.UnaryExpr, context *Context) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpr:
		val = EvaluateBinaryExpr(t, context)
	case *ast.UnaryExpr:
		val = EvaluateUnaryExpr(t, context)
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_LOGICAL_NOT:
		if bval, ok := val.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean())
		} else {
			panic(fmt.Errorf("BooleanValue interface is not support for %+v", val))
		}
	}
	return val
}

/*
EvaluateExpr calls EvaluateBinaryExpr. except EvaluateExpr
prevents calculate css slash as division.  otherwise it's the same as
EvaluateBinaryExpr.
*/
func EvaluateExpr(expr ast.Expr, context *Context) ast.Value {

	switch t := expr.(type) {

	case *ast.BinaryExpr:
		// For binary expression that is a CSS slash, we evaluate the expression as a literal string (unquoted)
		if t.IsCssSlash() {
			// return string object without quote
			return ast.NewString(0, t.Left.(*ast.Number).String()+"/"+t.Right.(*ast.Number).String(), nil)
		}
		return EvaluateBinaryExpr(t, context)

	case *ast.UnaryExpr:
		return EvaluateUnaryExpr(t, context)

	default:
		return ast.Value(expr)

	}

}

func EvaluateFunctionCall(fcall ast.FunctionCall, context *Context) ast.Value {
	if fun, ok := context.Functions.Get(fcall.Ident.Str); ok {

		_ = fun

	} else {
		panic("Function " + fcall.Ident.Str + " is undefined.")
	}
	return nil
}

/*
EvaluateBinaryExpr recursively.
*/
func EvaluateBinaryExpr(expr *ast.BinaryExpr, context *Context) ast.Value {
	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {

	case *ast.BinaryExpr:
		lval = EvaluateBinaryExpr(expr, context)

	case *ast.UnaryExpr:
		lval = EvaluateUnaryExpr(expr, context)

	case *ast.Variable:
		if varVal, ok := context.GetVariable(expr.Name); ok {
			lval = varVal.(ast.Expr)
		}

	default:
		lval = ast.Value(expr)
	}

	switch expr := expr.Right.(type) {

	case *ast.UnaryExpr:
		rval = EvaluateUnaryExpr(expr, context)

	case *ast.BinaryExpr:
		rval = EvaluateBinaryExpr(expr, context)

	case *ast.Variable:
		if varVal, ok := context.GetVariable(expr.Name); ok {
			rval = varVal.(ast.Expr)
		}

	default:
		rval = ast.Value(expr)
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpr(expr *ast.UnaryExpr, context *Context) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpr:
		val = EvaluateBinaryExpr(t, context)
	case *ast.UnaryExpr:
		val = EvaluateUnaryExpr(t, context)
	case *ast.Variable:
		if varVal, ok := context.GetVariable(t.Name); ok {
			val = varVal.(ast.Expr)
		}
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_NOP:
		// do nothing
	case ast.T_LOGICAL_NOT:
		if bVal, ok := val.(ast.BooleanValue); ok {
			val = ast.NewBoolean(bVal.Boolean())
		}
	case ast.T_MINUS:
		switch n := val.(type) {
		case *ast.Number:
			n.Value = -n.Value
		}
	}
	return val
}
