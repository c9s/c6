package runtime

import "fmt"
import "math"
import "c6/ast"

/*
Value
*/
type ComputeFunction func(a ast.Value, b ast.Value) ast.Value

const ValueTypeNum = 7

var computableMatrix [ValueTypeNum][ValueTypeNum]bool = [ValueTypeNum][ValueTypeNum]bool{
	/* NumberValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* LengthValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* HexColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBAColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},
}

/**
Each row: [5]ComputeFunction{ NumberValue, LengthValue, HexColorValue, RGBAColorValue, RGBColorValue }
*/
var computeFunctionMatrix [5][5]ComputeFunction = [5][5]ComputeFunction{

	/* NumberValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* LengthValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* HexColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBAColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},
}

func HexColorAddNumber(c *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := c.R + uint32(num.Value)
	g := c.G + uint32(num.Value)
	b := c.B + uint32(num.Value)
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{hex, r, g, b, nil}
}

func uintsub(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return 0
}

func HexColorSubNumber(c *ast.HexColor, num *ast.Number) *ast.HexColor {
	val := uint32(num.Value)
	r := uintsub(c.R, val)
	g := uintsub(c.G, val)
	b := uintsub(c.B, val)
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{hex, r, g, b, nil}
}

func HexColorMulNumber(color *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := uint32(math.Floor(float64(color.R) * num.Value))
	g := uint32(math.Floor(float64(color.G) * num.Value))
	b := uint32(math.Floor(float64(color.B) * num.Value))
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{hex, r, g, b, nil}
}

func HexColorDivNumber(color *ast.HexColor, num *ast.Number) *ast.HexColor {
	r := uint32(math.Floor(float64(color.R) / num.Value))
	g := uint32(math.Floor(float64(color.G) / num.Value))
	b := uint32(math.Floor(float64(color.B) / num.Value))
	hex := ast.Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	return &ast.HexColor{hex, r, g, b, nil}
}

/*
For color type, we treat them as vector. a vector can be computed with scalar or another vector.


Valid expression:

	#aaa / 3
	#bb * 3
	#bb - #cc
	3px - 3px
	3px + 3px

Invalid expression

- Number can't be the dividend.
- Hex Color can't divisor.

	3 / #aaa
	3 - #bbb
	3px - 3
	6em - 3px

*/

func NumberSubNumber(a *ast.Number, b *ast.Number) *ast.Number {
	var result = a.Value - b.Value
	return ast.NewNumber(result, nil)
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func NumberDivNumber(a *ast.Number, b *ast.Number) *ast.Number {
	var result = a.Value / b.Value
	return ast.NewNumber(result, nil)
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func NumberMulNumber(a *ast.Number, b *ast.Number) *ast.Number {
	var result = a.Value * b.Value
	return ast.NewNumber(result, nil)
}

func NumberAddNumber(a *ast.Number, b *ast.Number) *ast.Number {
	var result = a.Value + b.Value
	return ast.NewNumber(result, nil)
}

func LengthSubLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit != b.Unit {
		fmt.Printf("Incompatible unit %s != %s.  %v - %v \n", a.Unit, b.Unit, a, b)
		return nil
	}
	var result = a.Value - b.Value
	return ast.NewLength(result, a.Unit, nil)
}

func LengthAddLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit != b.Unit {
		fmt.Printf("Incompatible unit %s != %s.  %v + %v \n", a.Unit, b.Unit, a, b)
		return nil
	}
	var result = a.Value + b.Value
	return ast.NewLength(result, a.Unit, nil)
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func LengthDivLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit == ast.UNIT_NONE || b.Unit == ast.UNIT_NONE || a.Unit == b.Unit {
		var result = a.Value / b.Value
		var unit = ast.UNIT_NONE
		if a.Unit != ast.UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit != ast.UNIT_NONE {
			unit = b.Unit
		}
		return ast.NewLength(result, unit, nil)
	}
	return nil
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func LengthMulLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit == ast.UNIT_NONE || b.Unit == ast.UNIT_NONE || a.Unit == b.Unit {
		var result = a.Value * b.Value
		var unit = ast.UNIT_NONE
		if a.Unit != ast.UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit != ast.UNIT_NONE {
			unit = b.Unit
		}
		return ast.NewLength(result, unit, nil)
	}
	return nil
}

func LengthMulNumber(a *ast.Length, b *ast.Number) *ast.Length {
	return ast.NewLength(a.Value*b.Value, a.Unit, nil)
}

func NumberMulLength(a *ast.Number, b *ast.Length) *ast.Length {
	return ast.NewLength(a.Value*b.Value, b.Unit, nil)
}

func LengthDivNumber(a *ast.Length, b *ast.Number) *ast.Length {
	return ast.NewLength(a.Value/b.Value, a.Unit, nil)
}

func NumberDivLength(a *ast.Number, b *ast.Length) *ast.Length {
	panic("Number can't be divided by length")
	return nil
}

func RGBColorAddNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorSubNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorMulNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return ast.NewRGBColor(r, g, b, nil)
}

func RGBColorDivNumber(c *ast.RGBColor, n *ast.Number) *ast.RGBColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return ast.NewRGBColor(uint32(r), uint32(g), uint32(b), nil)
}

func RGBAColorAddNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = c.R + val
	var g = c.G + val
	var b = c.B + val
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorSubNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = uintsub(c.R, val)
	var g = uintsub(c.G, val)
	var b = uintsub(c.B, val)
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorMulNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = uint32(n.Value)
	var r = c.R * val
	var g = c.G * val
	var b = c.B * val
	return ast.NewRGBAColor(r, g, b, c.A, nil)
}

func RGBAColorDivNumber(c *ast.RGBAColor, n *ast.Number) *ast.RGBAColor {
	var val = n.Value
	var r = math.Floor(float64(c.R) / val)
	var g = math.Floor(float64(c.G) / val)
	var b = math.Floor(float64(c.B) / val)
	return ast.NewRGBAColor(uint32(r), uint32(g), uint32(b), c.A, nil)
}

func Compute(op ast.OpType, a ast.Value, b ast.Value) ast.Value {
	switch op {
	case ast.OpAdd:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberAddNumber(ta, tb)
			case *ast.HexColor:
				return HexColorAddNumber(tb, ta)
			}
		case *ast.Length:
			switch tb := b.(type) {
			case *ast.Length:
				return LengthAddLength(ta, tb)
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
	case ast.OpSub:
		switch ta := a.(type) {

		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberSubNumber(ta, tb)
			}

		case *ast.Length:
			switch tb := b.(type) {
			case *ast.Length:
				val := LengthSubLength(ta, tb)
				fmt.Printf("Substracted value: %+v\n", val)
				return val
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
	case ast.OpMul:
		switch ta := a.(type) {

		case *ast.Length:
			switch tb := b.(type) {
			case *ast.Length:
				return LengthMulLength(ta, tb)
			case *ast.Number:
				return LengthMulNumber(ta, tb)
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

func EvaluateBinaryExpression(expr *ast.BinaryExpression, symTable *SymTable) ast.Value {
	if expr.IsCssSlash() {
		// return string object without quote
		return ast.NewString(0, expr.Left.(*ast.Length).String()+"/"+expr.Right.(*ast.Length).String(), nil)
	}

	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {
	case *ast.BinaryExpression:
		lval = EvaluateBinaryExpression(expr, symTable)
	case *ast.UnaryExpression:
		lval = EvaluateUnaryExpression(expr, symTable)
	case *ast.Number, *ast.Length, *ast.HexColor:
		lval = ast.Value(expr)
	}
	switch expr := expr.Right.(type) {
	case *ast.UnaryExpression:
		rval = EvaluateUnaryExpression(expr, symTable)
	case *ast.BinaryExpression:
		rval = EvaluateBinaryExpression(expr, symTable)
	case *ast.Number, *ast.Length, *ast.HexColor:
		rval = ast.Value(expr)
	}
	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpression(expr *ast.UnaryExpression, symTable *SymTable) ast.Value {
	var val ast.Value = nil
	if bexpr, ok := expr.Expr.(*ast.BinaryExpression); ok {
		val = EvaluateBinaryExpression(bexpr, symTable)
	} else if uexpr, ok := expr.Expr.(*ast.UnaryExpression); ok {
		val = EvaluateUnaryExpression(uexpr, symTable)
	}

	// negative value
	if expr.Op == ast.OpSub {
		switch n := val.(type) {
		case *ast.Number:
			n.Value = -n.Value
		case *ast.Length:
			n.Value = -n.Value
		}
	}
	return val
}
