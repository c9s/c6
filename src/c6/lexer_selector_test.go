package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestLexerClassNameSelector(t *testing.T) {
	l := NewLexerWithString(`.class { }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameSelector(t *testing.T) {
	l := NewLexerWithString(`a {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameSelectorForDiv(t *testing.T) {
	l := NewLexerWithString(`div {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithUniversalSelector(t *testing.T) {
	l := NewLexerWithString(`* {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_UNIVERSAL_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`[href] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_BRACKET_RIGHT, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToUnquoteString(t *testing.T) {
	l := NewLexerWithString(`[lang=en] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_UNQUOTE_STRING, T_BRACKET_RIGHT, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToQQString(t *testing.T) {
	l := NewLexerWithString(`[lang="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_BRACKET_RIGHT, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorContainsQQString(t *testing.T) {
	l := NewLexerWithString(`[lang~="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_TILDE_EQUAL, T_QQ_STRING, T_BRACKET_RIGHT, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithAttributeSelectorAfterTagNameContainsQQString2(t *testing.T) {
	l := NewLexerWithString(`a[rel~="copyright"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_TILDE_EQUAL, T_QQ_STRING, T_BRACKET_RIGHT, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleSimpleSelectorGrouping(t *testing.T) {
	l := NewLexerWithString(`h1, h2, h3 { color: blue; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_TAGNAME_SELECTOR, T_COMMA, T_TAGNAME_SELECTOR, T_COMMA, T_TAGNAME_SELECTOR, T_BRACE_START,
		T_PROPERTY_NAME,
		T_COLON,
		T_CONSTANT,
		T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleAttributeSelectorGrouping(t *testing.T) {
	l := NewLexerWithString(`[type=text], [type=password], [type=checkbox] {}`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{

		T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_UNQUOTE_STRING, T_BRACKET_RIGHT, T_COMMA,
		T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_UNQUOTE_STRING, T_BRACKET_RIGHT, T_COMMA,
		T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_UNQUOTE_STRING, T_BRACKET_RIGHT,

		T_BRACE_START,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithCombinedAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`span[hello="Cleveland"][goodbye="Columbus"] { color: blue; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_TAGNAME_SELECTOR,
		T_AND_SELECTOR,
		T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_BRACKET_RIGHT,
		T_BRACKET_LEFT, T_ATTRIBUTE_NAME, T_EQUAL, T_QQ_STRING, T_BRACKET_RIGHT,
		T_BRACE_START,
		T_PROPERTY_NAME,
		T_COLON,
		T_CONSTANT,
		T_SEMICOLON,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithTagNameAndClassSelector(t *testing.T) {
	l := NewLexerWithString(`a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleUniversalSelectorPlusClassSelectorPlusAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts[href="http://google.com"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_UNIVERSAL_SELECTOR,
		T_AND_SELECTOR,
		T_CLASS_SELECTOR,
		T_AND_SELECTOR,
		T_BRACKET_LEFT,
		T_ATTRIBUTE_NAME,
		T_EQUAL,
		T_QQ_STRING,
		T_BRACKET_RIGHT,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleUniversalPlusClassSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_UNIVERSAL_SELECTOR,
		T_AND_SELECTOR,
		T_CLASS_SELECTOR,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleChildSelector(t *testing.T) {
	l := NewLexerWithString(`div.posts > a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR,
		T_GT,
		T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_CLASS_SELECTOR,
		T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithPseudoSelector(t *testing.T) {
	var testCases = []string{`:hover {  }`, `:link {  }`, `:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []TokenType{T_PSEUDO_SELECTOR, T_BRACE_START, T_BRACE_END})
		l.close()
	}
}

func TestLexerRuleWithTagNameAndPseudoSelector(t *testing.T) {
	var testCases = []string{`a:hover {  }`, `a:link {  }`, `a:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_PSEUDO_SELECTOR, T_BRACE_START, T_BRACE_END})
		l.close()
	}
}

func TestLexerRuleLangPseudoSelector(t *testing.T) {
	// html:lang(fr-ca) { quotes: '« ' ' »' }
	l := NewLexerWithString(`html:lang(fr-ca) {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_TAGNAME_SELECTOR, T_AND_SELECTOR, T_PSEUDO_SELECTOR, T_LANG_CODE, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithIdSelector(t *testing.T) {
	l := NewLexerWithString(`#myPost {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithIdSelectorWithDigits(t *testing.T) {
	l := NewLexerWithString(`#foo123 {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerPropertyEmValueMul(t *testing.T) {
	l := NewLexerWithString(`.foo { width: 1.3em * 10.2em }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_CLASS_SELECTOR, T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_FLOAT, T_UNIT_EM, T_MUL, T_FLOAT, T_UNIT_EM,
		T_BRACE_END})
	l.close()
}

func TestLexerPropertyPxValueMul(t *testing.T) {
	l := NewLexerWithString(`.foo { width: 1px * 3px }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_CLASS_SELECTOR, T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_INTEGER, T_UNIT_PX, T_MUL, T_INTEGER, T_UNIT_PX,
		T_BRACE_END})
	l.close()
}

func TestLexerRuleWithMultipleSelector(t *testing.T) {
	l := NewLexerWithString(`#foo123, .foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{T_ID_SELECTOR, T_COMMA, T_CLASS_SELECTOR, T_BRACE_START, T_BRACE_END})
	l.close()
}

func TestLexerRuleWithSubRuleWithParentSelector(t *testing.T) {
	l := NewLexerWithString(`.test { -webkit-transition: none;   &.foo { color: #fff; } }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []TokenType{
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_CONSTANT, T_SEMICOLON,
		T_PARENT_SELECTOR,
		T_AND_SELECTOR,
		T_CLASS_SELECTOR,
		T_BRACE_START,
		T_PROPERTY_NAME, T_COLON, T_HEX_COLOR, T_SEMICOLON,
		T_BRACE_END,
		T_BRACE_END})
	l.close()
}
