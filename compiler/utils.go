package compiler

import "strings"

func indent(level int) string {
	// two space
	return strings.Repeat("  ", level)
}
