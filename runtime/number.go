package runtime

import "github.com/c9s/c6/ast"
import "fmt"

func NumberComparable(a *ast.Number, b *ast.Number) bool {
	if a.Unit == nil && b.Unit == nil {
		return true
	}
	if a.Unit != nil && b.Unit != nil && a.Unit.Type == b.Unit.Type {
		return true
	}
	return true
}

func NumberSubNumber(a *ast.Number, b *ast.Number) *ast.Number {
	if a.Unit != nil && b.Unit != nil && a.Unit.Type != b.Unit.Type {
		fmt.Printf("Incompatible unit %s != %s.  %v - %v \n", a.Unit, b.Unit, a, b)
		return nil
	}
	var result = a.Value - b.Value
	return ast.NewNumber(result, a.Unit, nil)
}

func NumberAddNumber(a *ast.Number, b *ast.Number) *ast.Number {
	if (a.Unit == nil && b.Unit == nil) || (a.Unit != nil && b.Unit != nil && a.Unit.Type == b.Unit.Type) {
		return ast.NewNumber(a.Value+b.Value, a.Unit, nil)
	}
	fmt.Printf("Incompatible unit %s != %s.  %v + %v \n", a.Unit, b.Unit, a, b)
	return nil
}

/*
10px / 3, 10 / 3, 10px / 10px is allowed here
*/
func NumberDivNumber(a *ast.Number, b *ast.Number) *ast.Number {
	// for 10/2 and 10px/2px
	if (a.Unit == nil && b.Unit == nil) || (a.Unit != nil && b.Unit != nil && a.Unit.Type == b.Unit.Type) {
		return ast.NewNumber(a.Value/b.Value, nil, nil)
	}
	if a.Unit != nil && b.Unit == nil {
		return ast.NewNumber(a.Value/b.Value, a.Unit, nil)
	}
	return nil
}

/*
3 * 10px, 10px * 3, 10px * 10px is allowed here
*/
func NumberMulNumber(a *ast.Number, b *ast.Number) *ast.Number {
	if a.Unit == nil && b.Unit == nil {
		return ast.NewNumber(a.Value*b.Value, nil, nil)
	}

	if a.Unit == nil || b.Unit == nil || a.Unit.Type == b.Unit.Type {
		var result = a.Value * b.Value
		var unit *ast.Unit = nil
		if a.Unit != nil {
			unit = a.Unit
		}
		if b.Unit != nil {
			unit = b.Unit
		}
		return ast.NewNumber(result, unit, nil)
	}
	return nil
}
