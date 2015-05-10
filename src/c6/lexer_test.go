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
	AssertLexerTokenSequence(t, `@import url("test.css");`, []ast.TokenType{ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_SEMICOLON})
}

func TestLexerAtRuleImportWithUrlAndOneMediaType(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url("test.css") screen;`, []ast.TokenType{
		ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerAtRuleImportWithUrlAndTwoMediaType(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url("test.css") tv, projection;`, []ast.TokenType{
		ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_QQ_STRING, ast.T_PAREN_END, ast.T_IDENT, ast.T_COMMA, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerAtRuleImportWithUnquoteUrl(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url(http://foo.com/bar/test.css);`, []ast.TokenType{
		ast.T_IMPORT, ast.T_IDENT, ast.T_PAREN_START, ast.T_UNQUOTE_STRING, ast.T_PAREN_END, ast.T_SEMICOLON,
	})
}

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

func TestLexerVariableAssignmentWithCommentBlock(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: /* comment */ #fff /* comment 2 */;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON})
}

func TestLexerVariableAssignment(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON})
}

func TestLexerVariableAssignmentWithDefaultFlag(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff !default;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_DEFAULT, ast.T_SEMICOLON})
}

func TestLexerVariableAssignmentWithImportantFlag(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff !important;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_IMPORTANT, ast.T_SEMICOLON})
}

func TestLexerVariableAssignmentWithOptionalFlag(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #fff !optional;`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON, ast.T_HEX_COLOR, ast.T_OPTIONAL, ast.T_SEMICOLON})
}

func TestLexerVariableAssignmentWithInterp(t *testing.T) {
	AssertLexerTokenSequence(t, `$favorite: #{ 10 + 2 }px;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTERPOLATION_START, ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER, ast.T_INTERPOLATION_END, ast.T_LITERAL_CONCAT, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithIdent(t *testing.T) {
	AssertLexerTokenSequence(t, `$media: screen;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$feature: -webkit-min-device-pixel-ratio;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithFloatingNumber(t *testing.T) {
	AssertLexerTokenSequence(t, `$value: 1.5;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithPtValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 10pt;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PT, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithLengthValueStartWithDot(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: .3em;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithLengthValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 10px;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3em;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_EM, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3rem;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_REM, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3pt;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_PT, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3in;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_IN, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3cm;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_CM, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3ch;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_CH, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithDpiValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: 0.3dpi;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_DPI, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3dpcm;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_DPCM, ast.T_SEMICOLON,
	})
	AssertLexerTokenSequence(t, `$foo: 0.3dppx;`, []ast.TokenType{
		ast.T_VARIABLE, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_DPPX, ast.T_SEMICOLON,
	})
}

func TestLexerVariableAssignmentWithPercent(t *testing.T) {
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
		ast.T_LITERAL_CONCAT,
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
		ast.T_INTEGER, // 1
		ast.T_PLUS,
		ast.T_INTEGER, // 2
		ast.T_INTERPOLATION_END,
		ast.T_LITERAL_CONCAT, // px
		ast.T_IDENT,
		ast.T_INTEGER,
		ast.T_UNIT_PX,
		ast.T_SEMICOLON,
		ast.T_BRACE_END})
}

func TestLexerInterpolationComplex1(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: (10+10)#{ 2 * 2 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_PAREN_END,
		ast.T_LITERAL_CONCAT,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_MUL,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_LITERAL_CONCAT,
		ast.T_IDENT,
	})
}

func TestLexerMap(t *testing.T) {
	AssertLexerTokenSequence(t, `$var: (foo: 1, bar: 2);`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_COMMA,
		ast.T_IDENT, ast.T_COLON, ast.T_INTEGER,
		ast.T_PAREN_END,
	})
}

func TestLexerMapWithExtraComma(t *testing.T) {
	AssertLexerTokenSequence(t, `$var: (foo: 1, bar: 2, );`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_COMMA,
		ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_COMMA,
		ast.T_PAREN_END,
	})
}

func TestLexerMapWithExpression(t *testing.T) {
	AssertLexerTokenSequence(t, `$var: (foo: 2px + 3px, bar: $var2);`, []ast.TokenType{ast.T_VARIABLE, ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_PLUS, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_COMMA,
		ast.T_IDENT, ast.T_COLON, ast.T_VARIABLE,
		ast.T_PAREN_END,
	})
}

func TestLexerInterpolationComplex2(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: (10+10) #{ 2 * 2 } 10px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_PAREN_END,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_MUL,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_INTEGER,
		ast.T_UNIT_PX,
	})
}

func TestLexerInterpolationLeadingAndTrailing(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		padding: 10#{ 1 + 1 }#{ 2 + 2 }px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,

		ast.T_INTEGER,
		ast.T_LITERAL_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_LITERAL_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_LITERAL_CONCAT,

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

		ast.T_LITERAL_CONCAT,

		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END,

		ast.T_LITERAL_CONCAT,

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
		ast.T_LITERAL_CONCAT,
		ast.T_IDENT,
		ast.T_INTEGER,
		ast.T_UNIT_PX,
		ast.T_BRACE_END,
	})
}

func TestLexerCommentBlockBeforeRuleSet(t *testing.T) {
	AssertLexerTokenSequence(t, `
	/* comment block */
	.test .foo { }`, []ast.TokenType{
		ast.T_COMMENT_BLOCK,
		ast.T_CLASS_SELECTOR,
		ast.T_DESCENDANT_COMBINATOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_BRACE_END,
	})
}

func TestLexerCommentBlockBetweenSelectors(t *testing.T) {
	AssertLexerTokenSequence(t, `.test /* comment between selector and block */ .foo { }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_DESCENDANT_COMBINATOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_BRACE_END,
	})
}

func TestLexerCommentBlockBetweenSelectorAndBlock(t *testing.T) {
	AssertLexerTokenSequence(t, `.test /* comment between selector and block */ { }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_BRACE_END,
	})
}

func TestLexerCommentBlock(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		/* comment here */
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_COMMENT_BLOCK,
		ast.T_BRACE_END,
	})
}

/*
This is for microsoft filter functions
*/
func TestLexerMSFilterGradient(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		filter: progid:DXImageTransform.Microsoft.gradient(GradientType=0, startColorstr='#81a8cb', endColorstr='#4477a1');
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,

		ast.T_MS_PROGID,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_START,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_Q_STRING, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_Q_STRING,
		ast.T_PAREN_END,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerMSFilterBasicImage(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		filter: progid:DXImageTransform.Microsoft.BasicImage(rotation=2, mirror=1);
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,

		ast.T_MS_PROGID,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_START,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER,
		ast.T_PAREN_END,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerGradientFunction(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		background-image: -moz-linear-gradient(top, #81a8cb, #4477a1);
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_START,
		ast.T_IDENT,
		ast.T_COMMA,
		ast.T_HEX_COLOR,
		ast.T_COMMA,
		ast.T_HEX_COLOR,
		ast.T_PAREN_END,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerInterpolationPropertyNameSuffix(t *testing.T) {

	AssertLexerTokenSequence(t, `.test {
		padding-#{ 'left' }: 10px;
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_LITERAL_CONCAT,
		ast.T_INTERPOLATION_START,
		ast.T_Q_STRING,
		ast.T_INTERPOLATION_END,
		ast.T_COLON,
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
		ast.T_LITERAL_CONCAT,
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
		ast.T_LITERAL_CONCAT,
		ast.T_INTERPOLATION_START,
		ast.T_QQ_STRING,
		ast.T_INTERPOLATION_END,
		ast.T_LITERAL_CONCAT,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER,
		ast.T_PLUS,
		ast.T_INTEGER,
		ast.T_INTERPOLATION_END,
		ast.T_LITERAL_CONCAT,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func BenchmarkLexer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`.test, .foo, .bar { color: #fff; }`)
		var o = l.getOutput()
		l.run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.close()
	}
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

func TestLexerMapValue(t *testing.T) {
	AssertLexerTokenSequence(t, `$foo: (foo: 1, bar: 2);`, []ast.TokenType{
		ast.T_VARIABLE,
		ast.T_COLON,
		ast.T_PAREN_START,
		ast.T_IDENT,
		ast.T_COLON,
		ast.T_INTEGER,
		ast.T_COMMA,

		ast.T_IDENT,
		ast.T_COLON,
		ast.T_INTEGER,
		ast.T_PAREN_END,

		ast.T_SEMICOLON,
	})
}

/**********************************************************************
If Statement Test Case
***********************************************************************/
func TestLexerIfStatementTrueCondition(t *testing.T) {
	AssertLexerTokenSequence(t, `
	@if true {
		color: red;
	}
	`, []ast.TokenType{ast.T_IF, ast.T_TRUE, ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerIfStatementFalseCondition(t *testing.T) {
	AssertLexerTokenSequence(t, `
	@if false {
		color: red;
	}
	`, []ast.TokenType{ast.T_IF, ast.T_FALSE, ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerIfStatementTrueOrFalseCondition(t *testing.T) {
	AssertLexerTokenSequence(t, `
	@if true or false {
		color: red;
	}
	`, []ast.TokenType{
		ast.T_IF, ast.T_TRUE, ast.T_LOGICAL_OR, ast.T_FALSE, ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_END,
	})
}

func TestLexerIfStatementFunctionCallEqualToNumber(t *testing.T) {
	AssertLexerTokenSequence(t, `
	@if type-of(nth($x, 3)) == 10 {
	}
	`, []ast.TokenType{
		ast.T_IF, ast.T_FUNCTION_NAME, ast.T_PAREN_START, ast.T_FUNCTION_NAME, ast.T_PAREN_START, ast.T_VARIABLE, ast.T_COMMA, ast.T_INTEGER, ast.T_PAREN_END, ast.T_PAREN_END, ast.T_EQUAL, ast.T_INTEGER, ast.T_BRACE_START,
	})
}

/**********************************************************************

***********************************************************************/
func TestLexerForStatementSimpleFromThrough(t *testing.T) {
	AssertLexerTokenSequence(t, `@for $var from 1 through 20 {  }`, []ast.TokenType{
		ast.T_FOR, ast.T_VARIABLE, ast.T_FOR_FROM, ast.T_INTEGER, ast.T_FOR_THROUGH, ast.T_INTEGER, ast.T_BRACE_START, ast.T_BRACE_END,
	})
}

/**********************************************************************
Media Query Test Case
***********************************************************************/

func TestLexerMediaQueryCondition(t *testing.T) {
	AssertLexerTokenSequence(t, `@media screen and (orientation: landscape) { }`,
		[]ast.TokenType{ast.T_MEDIA, ast.T_IDENT, ast.T_LOGICAL_AND,
			ast.T_PAREN_START, ast.T_IDENT, ast.T_COLON, ast.T_IDENT, ast.T_PAREN_END,
			ast.T_BRACE_START,
			ast.T_BRACE_END,
		})
}

func TestLexerMediaQueryConditionWithExpressions(t *testing.T) {
	AssertLexerTokenSequence(t, `@media #{$media} and ($feature: $value) {
  .sidebar {
    width: 500px;
  }
}`,
		[]ast.TokenType{ast.T_MEDIA, ast.T_INTERPOLATION_START, ast.T_VARIABLE, ast.T_INTERPOLATION_END, ast.T_LOGICAL_AND,
			ast.T_PAREN_START, ast.T_VARIABLE, ast.T_COLON, ast.T_VARIABLE, ast.T_PAREN_END,
			ast.T_BRACE_START,
			ast.T_CLASS_SELECTOR,
			ast.T_BRACE_START, ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON, ast.T_BRACE_END,
			ast.T_BRACE_END,
		})
}

func TestLexerMediaQueryConditionSimpleMaxWidth(t *testing.T) {
	code := `
	@media (max-width: 1024px) {
		html, body, .container-base, .top-bar {
			width: 1024px;
		}
	}
	`
	AssertLexerTokenSequence(t, code, []ast.TokenType{
		ast.T_MEDIA, ast.T_PAREN_START, ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_PAREN_END, ast.T_BRACE_START,
		ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_CLASS_SELECTOR, ast.T_COMMA, ast.T_CLASS_SELECTOR, ast.T_BRACE_START,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON,
		ast.T_BRACE_END,
		ast.T_BRACE_END,
	})
}
