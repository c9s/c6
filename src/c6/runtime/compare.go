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

	case *ast.Number:
		switch b := bv.(type) {
		case *ast.Number:
			if (a.Unit == nil && b.Unit == nil) || (a.Unit != nil && b.Unit != nil) {
				return true
			}
			if a.Unit != nil && b.Unit != nil && a.Unit.Type != b.Unit.Type {
				return false
			}
			return true
		}
	}
	return false
}
