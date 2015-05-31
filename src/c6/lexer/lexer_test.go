package lexer

import "c6/ast"
import "testing"
import "github.com/stretchr/testify/assert"

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
	output := l.TokenStream()
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
	AssertLexerTokenSequence(t, `@import url("test.css");`, []ast.TokenType{ast.T_IMPORT, ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_QQ_STRING, ast.T_PAREN_CLOSE, ast.T_SEMICOLON})
}

func TestLexerAtRuleImportWithUrlAndOneMediaType(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url("test.css") screen;`, []ast.TokenType{
		ast.T_IMPORT, ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_QQ_STRING, ast.T_PAREN_CLOSE, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerAtRuleImportWithUrlAndTwoMediaType(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url("test.css") tv, projection;`, []ast.TokenType{
		ast.T_IMPORT, ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_QQ_STRING, ast.T_PAREN_CLOSE, ast.T_IDENT, ast.T_COMMA, ast.T_IDENT, ast.T_SEMICOLON,
	})
}

func TestLexerAtRuleImportWithUnquoteUrl(t *testing.T) {
	AssertLexerTokenSequence(t, `@import url(http://foo.com/bar/test.css);`, []ast.TokenType{
		ast.T_IMPORT, ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_UNQUOTE_STRING, ast.T_PAREN_CLOSE, ast.T_SEMICOLON,
	})
}

func TestLexerRuleWithOneProperty(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { color: #fff; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithTwoProperty(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { color: #fff; background: #fff; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithPropertyValueComma(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { font-family: Arial, sans-serif }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_COMMA, ast.T_IDENT,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerRuleWithVendorPrefixPropertyName(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { -webkit-transition: none; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithVariableAsPropertyValue(t *testing.T) {
	AssertLexerTokenSequence(t, `.test { color: $favorite; }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_VARIABLE, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
}

func TestLexerCommentBlockBeforeRuleSet(t *testing.T) {
	AssertLexerTokenSequence(t, `
	/* comment block */
	.test .foo { }`, []ast.TokenType{
		ast.T_COMMENT_BLOCK,
		ast.T_CLASS_SELECTOR,
		ast.T_DESCENDANT_COMBINATOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerCommentBlockBetweenSelectors(t *testing.T) {
	AssertLexerTokenSequence(t, `.test /* comment between selector and block */ .foo { }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_DESCENDANT_COMBINATOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerCommentBlockBetweenSelectorAndBlock(t *testing.T) {
	AssertLexerTokenSequence(t, `.test /* comment between selector and block */ { }`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerCommentBlock(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		/* comment here */
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_COMMENT_BLOCK,
		ast.T_BRACE_CLOSE,
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
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,

		ast.T_MS_PROGID,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_OPEN,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_Q_STRING, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_Q_STRING,
		ast.T_PAREN_CLOSE,
		ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerMSFilterBasicImage(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		filter: progid:DXImageTransform.Microsoft.BasicImage(rotation=2, mirror=1);
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,

		ast.T_MS_PROGID,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_OPEN,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER, ast.T_COMMA,
		ast.T_MS_PARAM_NAME, ast.T_ATTR_EQUAL, ast.T_INTEGER,
		ast.T_PAREN_CLOSE,
		ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerGradientFunction(t *testing.T) {
	AssertLexerTokenSequence(t, `.test {
		background-image: -moz-linear-gradient(top, #81a8cb, #4477a1);
	}`, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_FUNCTION_NAME,
		ast.T_PAREN_OPEN,
		ast.T_IDENT,
		ast.T_COMMA,
		ast.T_HEX_COLOR,
		ast.T_COMMA,
		ast.T_HEX_COLOR,
		ast.T_PAREN_CLOSE,
		ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerPagedMedia(t *testing.T) {
	AssertLexerTokenSequence(t, `@page { margin: 2cm }`,
		[]ast.TokenType{ast.T_PAGE, ast.T_BRACE_OPEN, ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_CM, ast.T_BRACE_CLOSE})
}

func TestLexerPagedMediaWithPseudoSelector(t *testing.T) {
	AssertLexerTokenSequence(t, `@page :left { margin: 2cm }`,
		[]ast.TokenType{ast.T_PAGE, ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_CM, ast.T_BRACE_CLOSE})
}

/**********************************************************************
Media Query Test Case
***********************************************************************/

func TestLexerMediaQueryCondition(t *testing.T) {
	AssertLexerTokenSequence(t, `@media screen and (orientation: landscape) { }`,
		[]ast.TokenType{ast.T_MEDIA, ast.T_IDENT, ast.T_LOGICAL_AND,
			ast.T_PAREN_OPEN, ast.T_IDENT, ast.T_COLON, ast.T_IDENT, ast.T_PAREN_CLOSE,
			ast.T_BRACE_OPEN,
			ast.T_BRACE_CLOSE,
		})
}

func TestLexerMediaQueryConditionWithExprs(t *testing.T) {
	AssertLexerTokenSequence(t, `@media #{$media} and ($feature: $value) {
  .sidebar {
    width: 500px;
  }
}`,
		[]ast.TokenType{ast.T_MEDIA, ast.T_INTERPOLATION_START, ast.T_VARIABLE, ast.T_INTERPOLATION_END, ast.T_LOGICAL_AND,
			ast.T_PAREN_OPEN, ast.T_VARIABLE, ast.T_COLON, ast.T_VARIABLE, ast.T_PAREN_CLOSE,
			ast.T_BRACE_OPEN,
			ast.T_CLASS_SELECTOR,
			ast.T_BRACE_OPEN, ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON, ast.T_BRACE_CLOSE,
			ast.T_BRACE_CLOSE,
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
		ast.T_MEDIA, ast.T_PAREN_OPEN, ast.T_IDENT, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_PAREN_CLOSE, ast.T_BRACE_OPEN,
		ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_CLASS_SELECTOR, ast.T_COMMA, ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE,
		ast.T_BRACE_CLOSE,
	})
}

func TestLexerFontFace(t *testing.T) {
	var code = `
@font-face {
  font-family: MyGentium;
  src: local(Gentium Bold),    /* full font name */
       local(Gentium-Bold),    /* Postscript name */
       url(GentiumBold.woff);  /* otherwise, download it */
  font-weight: bold;
}
`
	AssertLexerTokenSequence(t, code, []ast.TokenType{
		ast.T_FONT_FACE, ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,

		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON,
		ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_UNQUOTE_STRING, ast.T_PAREN_CLOSE, ast.T_COMMA,
		ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_UNQUOTE_STRING, ast.T_PAREN_CLOSE, ast.T_COMMA,
		ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN, ast.T_UNQUOTE_STRING, ast.T_PAREN_CLOSE, ast.T_SEMICOLON,

		ast.T_COMMENT_BLOCK,

		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,

		ast.T_BRACE_CLOSE,
	})
}

func BenchmarkLexerBasic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`.test, .foo, .bar { color: #fff; }`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerTypeSelector(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`div { }`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerIdSelector(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`#myId { }`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerClassSelector(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`.foo-bar-zoo { }`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerAttributeSelector(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`
		[href*=ftp] { }
		`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerImportRule(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`@import url(../test.scss);`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}

func BenchmarkLexerVariableDeclaration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Fib(10)
		var l = NewLexerWithString(`$color: (3+1) * 4;`)
		var o = l.TokenStream()
		l.Run()
		var token = <-o
		for ; token != nil; token = <-o {
		}
		l.Close()
	}
}
