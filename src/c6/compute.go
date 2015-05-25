package c6

import "github.com/c9s/c6/src/c6/ast"
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

func IsConstantExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		if IsConstantValue(e.Left) && IsConstantValue(e.Right) {
			return true
		}
	case *ast.UnaryExpression:
		if IsConstantValue(e.Expr) {
			return true
		}
	}
	return false
}

func IsConstantValue(val ast.Value) bool {
	switch val.(type) {
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor, *ast.HSLColor, *ast.HSVColor, *ast.Boolean:
		return true
	}
	return false
}

func EvaluateExpressionInBooleanContext(anyexpr ast.Expression, context *Context) ast.Value {
	switch expr := anyexpr.(type) {

	case *ast.BinaryExpression:
		return EvaluateBinaryExpressionInBooleanContext(expr, context)

	case *ast.UnaryExpression:
		return EvaluateUnaryExpressionInBooleanContext(expr, context)

	default:
		if bval, ok := expr.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean())
		}
	}
	return nil
}

func EvaluateBinaryExpressionInBooleanContext(expr *ast.BinaryExpression, context *Context) ast.Value {

	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {
	case *ast.UnaryExpression:
		lval = EvaluateUnaryExpressionInBooleanContext(expr, context)

	case *ast.BinaryExpression:
		lval = EvaluateBinaryExpressionInBooleanContext(expr, context)

	default:
		lval = expr
	}

	switch expr := expr.Right.(type) {
	case *ast.UnaryExpression:
		rval = EvaluateUnaryExpressionInBooleanContext(expr, context)

	case *ast.BinaryExpression:
		rval = EvaluateBinaryExpressionInBooleanContext(expr, context)

	default:
		rval = expr
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpressionInBooleanContext(expr *ast.UnaryExpression, context *Context) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpression:
		val = EvaluateBinaryExpression(t, context)
	case *ast.UnaryExpression:
		val = EvaluateUnaryExpression(t, context)
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
EvaluateExpression calls EvaluateBinaryExpression. except EvaluateExpression
prevents calculate css slash as division.  otherwise it's the same as
EvaluateBinaryExpression.
*/
func EvaluateExpression(expr ast.Expression, context *Context) ast.Value {

	switch t := expr.(type) {

	case *ast.BinaryExpression:
		// For binary expression that is a CSS slash, we evaluate the expression as a literal string (unquoted)
		if t.IsCssSlash() {
			// return string object without quote
			return ast.NewString(0, t.Left.(*ast.Number).String()+"/"+t.Right.(*ast.Number).String(), nil)
		}
		return EvaluateBinaryExpression(t, context)

	case *ast.UnaryExpression:
		return EvaluateUnaryExpression(t, context)

	default:
		return ast.Value(expr)

	}

	// shouldn't call here.
	if IsConstantValue(expr) {
		return ast.Value(expr)
	}

	panic("EvaluateExpression: Unsupported expression type")
	return nil
}

func EvaluateFunctionCall(fcall ast.FunctionCall, context *Context) ast.Value {
	if fun, ok := context.Functions.Get(fcall.Ident); ok {

		_ = fun

	} else {
		panic("Function " + fcall.Ident + " is undefined.")
	}
	return nil
}

/*
EvaluateBinaryExpression recursively.
*/
func EvaluateBinaryExpression(expr *ast.BinaryExpression, context *Context) ast.Value {
	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {

	case *ast.BinaryExpression:
		lval = EvaluateBinaryExpression(expr, context)

	case *ast.UnaryExpression:
		lval = EvaluateUnaryExpression(expr, context)

	case *ast.Variable:
		if varVal, ok := context.GetVariable(expr.Name); ok {
			lval = varVal.(ast.Expression)
		}

	default:
		lval = ast.Value(expr)
	}

	switch expr := expr.Right.(type) {

	case *ast.UnaryExpression:
		rval = EvaluateUnaryExpression(expr, context)

	case *ast.BinaryExpression:
		rval = EvaluateBinaryExpression(expr, context)

	case *ast.Variable:
		if varVal, ok := context.GetVariable(expr.Name); ok {
			rval = varVal.(ast.Expression)
		}

	default:
		rval = ast.Value(expr)
	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpression(expr *ast.UnaryExpression, context *Context) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpression:
		val = EvaluateBinaryExpression(t, context)
	case *ast.UnaryExpression:
		val = EvaluateUnaryExpression(t, context)
	case *ast.Variable:
		if varVal, ok := context.GetVariable(t.Name); ok {
			val = varVal.(ast.Expression)
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
