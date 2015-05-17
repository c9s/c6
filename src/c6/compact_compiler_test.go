package c6

import "testing"
import _ "github.com/stretchr/testify/assert"

func AssertCompile(t *testing.T, code string, expected string) {
	var context = NewContext()
	var parser = NewParser(context)
	var stmts = parser.ParseScss(code)

	var compiler = NewCompactCompiler(context)
	var out = compiler.Compile(stmts)
	_ = out
	// assert.Equal(t, expected, out)
}

func TestCompileSimple(t *testing.T) {
	AssertCompile(t,
		`* { }`,
		`* { }`)
}
