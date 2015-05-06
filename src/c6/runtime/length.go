package runtime

import "c6/ast"
import "fmt"

func LengthSubLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit.Type != b.Unit.Type {
		fmt.Printf("Incompatible unit %s != %s.  %v - %v \n", a.Unit, b.Unit, a, b)
		return nil
	}
	var result = a.Value - b.Value
	return ast.NewLength(result, a.Unit, nil)
}

func LengthAddLength(a *ast.Length, b *ast.Length) *ast.Length {
	if a.Unit.Type != b.Unit.Type {
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
	if a.Unit.Type == ast.T_UNIT_NONE || b.Unit.Type == ast.T_UNIT_NONE || a.Unit.Type == b.Unit.Type {
		var result = a.Value / b.Value
		var unit *ast.Unit = nil
		if a.Unit.Type != ast.T_UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit.Type != ast.T_UNIT_NONE {
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
	if a.Unit.Type == ast.T_UNIT_NONE || b.Unit.Type == ast.T_UNIT_NONE || a.Unit.Type == b.Unit.Type {
		var result = a.Value * b.Value
		var unit *ast.Unit = nil
		if a.Unit.Type != ast.T_UNIT_NONE {
			unit = a.Unit
		}
		if b.Unit.Type != ast.T_UNIT_NONE {
			unit = b.Unit
		}
		return ast.NewLength(result, unit, nil)
	}
	return nil
}

func LengthMulNumber(a *ast.Length, b *ast.Number) *ast.Length {
	return ast.NewLength(a.Value*b.Value, a.Unit, nil)
}

func LengthDivNumber(a *ast.Length, b *ast.Number) *ast.Length {
	return ast.NewLength(a.Value/b.Value, a.Unit, nil)
}
