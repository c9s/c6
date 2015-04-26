package c6

import "c6/ast"
import "testing"
import "github.com/stretchr/testify/assert"

// not used right now.
type TestCase struct {
	Sample string
	Result bool
}

func TestLexerNext(t *testing.T) {
	l := NewLexerWithString(`.test {  }`)
	assert.NotNil(t, l)

	var r rune
	r = l.next()
	assert.Equal(t, '.', r)

	r = l.next()
	assert.Equal(t, 't', r)

	r = l.next()
	assert.Equal(t, 'e', r)

	r = l.next()
	assert.Equal(t, 's', r)

	r = l.next()
	assert.Equal(t, 't', r)
}

func TestLexerMatch(t *testing.T) {
	l := NewLexerWithString(`.foo {  }`)
	assert.NotNil(t, l)
	assert.False(t, l.match(".bar"))
	assert.True(t, l.match(".foo"))
}

func TestLexerAccept(t *testing.T) {
	l := NewLexerWithString(`.foo {  }`)
	assert.NotNil(t, l)
	assert.True(t, l.accept("."))
	assert.True(t, l.accept("f"))
	assert.True(t, l.accept("o"))
	assert.True(t, l.accept("o"))
	assert.True(t, l.accept(" "))
	assert.True(t, l.accept("{"))
}

func TestLexerIgnoreSpace(t *testing.T) {
	l := NewLexerWithString(`       .test {  }`)
	assert.NotNil(t, l)

	l.ignoreSpaces()

	var r rune
	r = l.next()
	assert.Equal(t, '.', r)

	l.backup()
	assert.True(t, l.match(".test"))
}

func TestLexerString(t *testing.T) {
	l := NewLexerWithString(`   "foo"`)
	output := l.getOutput()
	assert.NotNil(t, l)
	l.til("\"")
	lexString(l)
	token := <-output
	assert.Equal(t, ast.T_QQ_STRING, token.Type)
}

func TestLexerTil(t *testing.T) {
	l := NewLexerWithString(`"foo"`)
	assert.NotNil(t, l)
	l.til("\"")
	assert.Equal(t, 0, l.Offset)
	l.next() // skip the quote

	l.til("\"")
	assert.Equal(t, 4, l.Offset)
}

func TestLexerAtRuleImport(t *testing.T) {
	l := NewLexerWithString(`@import "test.css";`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_IMPORT, ast.T_QQ_STRING, ast.T_SEMICOLON})
	l.close()
}

func TestLexerAtRuleImportWithUrl(t *testing.T) {
	l := NewLexerWithString(`@import url("test.css");`)
	assert.NotNil(t, l)
	l.run()
	tokens := AssertTokenSequence(t, l, []ast.TokenType{ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_SEMICOLON})

	for _, tok := range tokens {
		if tok.Type == ast.T_QQ_STRING {
			assert.Equal(t, `test.css`, tok.Str)
		}
	}

	l.close()
}

func TestLexerAtRuleImportWithUrlAndMediaList(t *testing.T) {
	l := NewLexerWithString(`@import url("test.css") screen;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_MEDIA, ast.T_SEMICOLON})
	l.close()
}

func TestLexerAtRuleImportWithUnquoteUrl(t *testing.T) {
	l := NewLexerWithString(`@import url(http://foo.com/bar/test.css);`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_UNQUOTE_STRING, ast.T_PAREN_END, ast.T_SEMICOLON})
	l.close()
}

/*
func TestLexerAtRuleImportWithQuoteUrl(t *testing.T) {
	l := NewLexerWithString(`@import url("test.css");`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_IMPORT, ast.T_QQ_STRING, ast.T_SEMICOLON})
	l.close()
}
*/

func TestLexerRuleWithOneProperty(t *testing.T) {
	l := NewLexerWithString(`.test { color: #fff; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTwoProperty(t *testing.T) {
	l := NewLexerWithString(`.test { color: #fff; background: #fff; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerRuleWithPropertyValueComma(t *testing.T) {
	l := NewLexerWithString(`.test { font-family: Arial, sans-serif }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_IDENT, ast.T_COMMA, ast.T_IDENT,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerRuleWithVendorPrefixPropertyName(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: none; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerRuleWithVariableAsPropertyValue(t *testing.T) {
	l := NewLexerWithString(`.test { color: $favorite; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_VARIABLE, ast.T_SEMICOLON,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerVariableAssignment(t *testing.T) {
	l := NewLexerWithString(`$favorite: #fff;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON})
	l.close()
}

func TestLexerVariableWithPtValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 10pt;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PT, ast.T_SEMICOLON,
	})
	l.close()
}

func TestLexerVariableWithPxValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 10px;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON,
	})
	l.close()
}

func TestLexerVariableWithEmValue(t *testing.T) {
	l := NewLexerWithString(`$foo: 0.3em;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
	l.close()
}

func TestLexerVariableWithPercent(t *testing.T) {
	l := NewLexerWithString(`width: 20%;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PERCENT, ast.T_SEMICOLON,
	})
	l.close()
}

func TestLexerMultipleVariableAssignment(t *testing.T) {
	l := NewLexerWithString(`$favorite: #fff; $foo: 10em;`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
	l.close()
}

func TestLexerInterpolationPropertyValue(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: #{ 1 + 2 }px }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_BRACE_END})
	l.close()
}

func TestLexerInterpolationPropertyName(t *testing.T) {

}

func TestLexerRuleWithSubRule(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: none;   .foo { color: #fff; } }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END,
		ast.T_BRACE_END})
	l.close()
}
