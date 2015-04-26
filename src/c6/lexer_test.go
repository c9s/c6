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
	AssertLexerTokenSequence(t, `@import "test.css";`, []ast.TokenType{ast.T_IMPORT, ast.T_QQ_STRING, ast.T_SEMICOLON})
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
	AssertLexerTokenSequence(t, `@import url("test.css") screen;`, []ast.TokenType{
		ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_MEDIA, ast.T_SEMICOLON,
	})
}

func TestLexerAtRuleImportWithUnquoteUrl(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url(http://foo.com/bar/test.css);`, []ast.TokenType{
		ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_UNQUOTE_STRING, ast.T_PAREN_END, ast.T_SEMICOLON,
	})
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
	AssertLexerTokenSequence(t, `.test { color: #fff; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerRuleWithTwoProperty(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { color: #fff; background: #fff; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerRuleWithPropertyValueComma(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { font-family: Arial, sans-serif }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_COMMA, ast.T_IDENT,
		ast.T_BRACE_END,
	})
}

func TestLexerRuleWithVendorPrefixPropertyName(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { -webkit-transition: none; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerRuleWithVariableAsPropertyValue(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { color: $favorite; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_VARIABLE, ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerVariableAssignment(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON})
}

func TestLexerVariableWithPtValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 10pt;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PT, ast.T_SEMICOLON,
	})
}

func TestLexerVariableWithPxValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 10px;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON,
	})
}

func TestLexerVariableWithEmValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 0.3em;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
}

func TestLexerVariableWithPercent(t *testing.T) {
	AssertLexerTokenSequence(t, `width: 20%;`, []ast.TokenType{
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PERCENT, ast.T_SEMICOLON,
	})
}

func TestLexerMultipleVariableAssignment(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff; $foo: 10em;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
}

func TestLexerInterpolationPropertyValue(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		-webkit-transition: #{ 1 + 2 }px
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_BRACE_END})
}

func TestLexerInterpolationPropertyValueList(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: #{ 1 + 2 }px 10px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_INTEGER,
		ast.T_UNIT_PX,
		ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerInterpolationLeadingAndTrailing(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: 10#{ 1 + 1 }#{ 2 + 2 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,

		ast.T_INTEGER,
		ast.T_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_CONCAT,

		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerInterpolationConcatInterpolation(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: #{ 1 + 2 }#{ 3 + 4 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_CONCAT,

		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerInterpolationPropertyValueListWithoutSemiColon(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: #{ 1 + 2 }px 10px
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_INTEGER,
		ast.T_UNIT_PX,
		ast.T_BRACE_END,
	})
}

func TestLexerInterpolationPropertyName(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		#{ "foo" }: #{ 1 + 2 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_INTERPOLATION_START,
		ast.T_QQ_STRING,
		ast.T_INTERPOLATION_END,
		ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerInterpolationPropertyName2(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		-#{ "moz" }-border-radius: #{ 1 + 2 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_CONCAT,
		ast.T_INTERPOLATION_START,
		ast.T_QQ_STRING,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_CONCAT,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerRuleWithSubRule(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		-webkit-transition: none;
		.foo { color: #fff; }
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_END,
		ast.T_BRACE_END,
	})
}
