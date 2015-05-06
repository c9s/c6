package runtime

import "c6/ast"

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

func NumberMulLength(a *ast.Number, b *ast.Length) *ast.Length {
	return ast.NewLength(a.Value*b.Value, b.Unit, nil)
}

func NumberDivLength(a *ast.Number, b *ast.Length) *ast.Length {
	panic("Number can't be divided by length")
	return nil
}
