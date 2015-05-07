package runtime

import "c6/ast"

func IsComparable(av ast.Value, bv ast.Value) bool {
	switch a := av.(type) {

	case *ast.HexColor:
		switch bv.(type) {
		case *ast.HexColor:
			return true
		}

	case *ast.RGBColor:
		switch bv.(type) {
		case *ast.RGBColor:
			return true
		}

	case *ast.RGBAColor:
		switch bv.(type) {
		case *ast.RGBAColor:
			return true
		}

	case *ast.Boolean:
		switch bv.(type) {
		case *ast.Boolean:
			return true
		}

	/*
		For numbers, we compare the unit
	*/
	case *ast.Number:
		switch b := bv.(type) {
		case *ast.Number:
			return NumberComparable(a, b)
		}
	}
	return false
}
