package c6

import "testing"
import "github.com/stretchr/testify/assert"

func AssertCompile(t *testing.T, code string, expected string) {
	var context = NewContext()
	var parser = NewParser(context)
	var stmts = parser.ParseScss(code)
	var compiler = NewCompactCompiler(context)
	var out = compiler.Compile(stmts)
	_ = out
	// assert.Equal(t, expected, out)
}

func TestCompilerCompliant(t *testing.T) {
	var context = NewContext()
	var compiler = NewCompactCompiler(context)
	compiler.EnableCompliant(CSS3Compliant)
	compiler.EnableCompliant(IE7Compliant)
	assert.True(t, compiler.HasCompliant(CSS3Compliant))
	assert.False(t, compiler.HasCompliant(CSS4Compliant))

	assert.True(t, compiler.HasCompliant(IE7Compliant))
	assert.False(t, compiler.HasCompliant(IE8Compliant))
}

func TestCompileSimple(t *testing.T) {
	AssertCompile(t,
		`* { }`,
		`* { }`)
}
