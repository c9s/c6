package compiler

import "testing"
import "c6/runtime"
import "c6/parser"
import "github.com/stretchr/testify/assert"

func AssertCompile(t *testing.T, code string, expected string) {
	var context = runtime.NewContext()
	var parser = parser.NewParser(context)
	var stmts = parser.ParseScss(code)
	var compiler = NewCompactCompiler(context)
	var out = compiler.CompileString(stmts)
	assert.Equal(t, expected, out)
}

func TestCompilerCompliant(t *testing.T) {
	var context = runtime.NewContext()
	var compiler = NewCompactCompiler(context)
	compiler.EnableCompliant(CSS3Compliant)
	compiler.EnableCompliant(IE7Compliant)
	assert.True(t, compiler.HasCompliant(CSS3Compliant))
	assert.False(t, compiler.HasCompliant(CSS4Compliant))

	assert.True(t, compiler.HasCompliant(IE7Compliant))
	assert.False(t, compiler.HasCompliant(IE8Compliant))
}

func TestCompileUniversalSelector(t *testing.T) {
	AssertCompile(t,
		`* { }`,
		`* {}`)
}

func TestCompileClassSelector(t *testing.T) {
	AssertCompile(t,
		`.foo-bar { }`,
		`.foo-bar {}`)
}

func TestCompileIdSelector(t *testing.T) {
	AssertCompile(t,
		`#myId { }`,
		`#myId {}`)
}

func TestCompileAttributeSelector(t *testing.T) {
	AssertCompile(t,
		`[type=text] { }`,
		`[type=text] {}`)
}

func TestCompileAttributeSelectorWithTypeName(t *testing.T) {
	AssertCompile(t,
		`input[type=text] { }`,
		`input[type=text] {}`)
}

func TestCompileSelectorGroup(t *testing.T) {
	AssertCompile(t,
		`html, span, div { }`,
		`html, span, div {}`)
}

func TestCompileCompoundSelector1(t *testing.T) {
	AssertCompile(t,
		`*.foo.bar { }`,
		`*.foo.bar {}`)
}

func TestCompileCompoundSelector2(t *testing.T) {
	AssertCompile(t,
		`div.foo.bar[href$=pdf] { }`,
		`div.foo.bar[href$=pdf] {}`)
}

func TestCompileComplexSelector(t *testing.T) {
	AssertCompile(t,
		`*.foo.bar > .posts { }`,
		`*.foo.bar > .posts {}`)
}
