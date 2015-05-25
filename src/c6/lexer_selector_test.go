package c6

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/c9s/c6/src/c6/ast"

func TestLexerClassNameSelector(t *testing.T) {
	l := NewLexerWithString(`.class { }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithTagNameSelector(t *testing.T) {
	l := NewLexerWithString(`a {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithTagNameSelectorWithProperty(t *testing.T) {
	AssertLexerTokenSequence(t, `div { width: 200px; }`, []ast.TokenType{
		ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_SEMICOLON, ast.T_BRACE_CLOSE,
	})
}

func TestLexerRuleWithTagNameSelectorForDiv(t *testing.T) {
	l := NewLexerWithString(`div {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithUniversalSelector(t *testing.T) {
	l := NewLexerWithString(`* {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_UNIVERSAL_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`[href] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_BRACKET_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToUnquoteString(t *testing.T) {
	l := NewLexerWithString(`[lang=en] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_UNQUOTE_STRING, ast.T_BRACKET_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithAttributeSelectorEqualToQQString(t *testing.T) {
	l := NewLexerWithString(`[lang="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_QQ_STRING, ast.T_BRACKET_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithAttributeSelectorContainsQQString(t *testing.T) {
	l := NewLexerWithString(`[lang~="en"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_INCLUDE_MATCH, ast.T_QQ_STRING, ast.T_BRACKET_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithAttributeSelectorAfterTagNameContainsQQString2(t *testing.T) {
	l := NewLexerWithString(`a[rel~="copyright"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_INCLUDE_MATCH, ast.T_QQ_STRING, ast.T_BRACKET_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleSimpleSelectorGrouping(t *testing.T) {
	l := NewLexerWithString(`h1, h2, h3 { color: blue; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleAttributeSelectorGrouping(t *testing.T) {
	l := NewLexerWithString(`[type=text], [type=password], [type=checkbox] {}`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{

		ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_UNQUOTE_STRING, ast.T_BRACKET_CLOSE, ast.T_COMMA,
		ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_UNQUOTE_STRING, ast.T_BRACKET_CLOSE, ast.T_COMMA,
		ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_UNQUOTE_STRING, ast.T_BRACKET_CLOSE,

		ast.T_BRACE_OPEN,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithCombinedAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`span[hello="Cleveland"][goodbye="Columbus"] { color: blue; }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_TYPE_SELECTOR,
		ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_QQ_STRING, ast.T_BRACKET_CLOSE,
		ast.T_BRACKET_OPEN, ast.T_ATTRIBUTE_NAME, ast.T_ATTR_EQUAL, ast.T_QQ_STRING, ast.T_BRACKET_CLOSE,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN,
		ast.T_COLON,
		ast.T_IDENT,
		ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithTagNameAndClassSelector(t *testing.T) {
	l := NewLexerWithString(`a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleForDescendantTagNameSelectorWithoutSpace(t *testing.T) {
	l := NewLexerWithString(`div input{}`)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_DESCENDANT_COMBINATOR, ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleForDescendantTagNameSelector(t *testing.T) {
	l := NewLexerWithString(`div input {  }`)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_DESCENDANT_COMBINATOR, ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleForDescendantClassSelector(t *testing.T) {
	l := NewLexerWithString(`.foo .bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_CLASS_SELECTOR, ast.T_DESCENDANT_COMBINATOR, ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleForDescendantClassSelectorAndTagNameSelector(t *testing.T) {
	l := NewLexerWithString(`div.foo span.bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_CLASS_SELECTOR, ast.T_DESCENDANT_COMBINATOR, ast.T_TYPE_SELECTOR, ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleAttributeSelectorWithInterpolationInAttributeName(t *testing.T) {
	l := NewLexerWithString(`[#{ $foo }] {  }`)
	assert.NotNil(t, l)
	l.run()
	output := l.getOutput()
	var token = <-output
	token = <-output
	assert.True(t, token.ContainsInterpolation)
	close(output)
}

func TestLexerRuleAttributeSelectorWithInterpolationInAttributeNameInTheMiddle(t *testing.T) {
	l := NewLexerWithString(`[data-#{ $foo }-type] {  }`)
	assert.NotNil(t, l)
	l.run()
	output := l.getOutput()
	var token = <-output
	token = <-output
	assert.True(t, token.ContainsInterpolation)
	close(output)
}

func TestLexerRuleAttributeSelectorWithInterpolationInAttributeName2(t *testing.T) {
	l := NewLexerWithString(`[#{ $foo }="http://google.com"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_BRACKET_OPEN,
		ast.T_ATTRIBUTE_NAME,
		ast.T_ATTR_EQUAL,
		ast.T_QQ_STRING,
		ast.T_BRACKET_CLOSE,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleUniversalSelectorPlusClassSelectorPlusAttributeSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts[href="http://google.com"] {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_UNIVERSAL_SELECTOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACKET_OPEN,
		ast.T_ATTRIBUTE_NAME,
		ast.T_ATTR_EQUAL,
		ast.T_QQ_STRING,
		ast.T_BRACKET_CLOSE,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleUniversalPlusClassSelector(t *testing.T) {
	l := NewLexerWithString(`*.posts {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_UNIVERSAL_SELECTOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleAdjacentCombinator(t *testing.T) {
	l := NewLexerWithString(`.cover + .content {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_ADJACENT_SIBLING_COMBINATOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleChildCombinator(t *testing.T) {
	l := NewLexerWithString(`div.posts > a.foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_TYPE_SELECTOR, ast.T_CLASS_SELECTOR,
		ast.T_CHILD_COMBINATOR,
		ast.T_TYPE_SELECTOR, ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithPseudoSelector(t *testing.T) {
	var testCases = []string{`:hover {  }`, `:link {  }`, `:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []ast.TokenType{ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
		l.close()
	}
}

func TestLexerRuleWithTagNameAndPseudoSelector(t *testing.T) {
	var testCases = []string{`a:hover {  }`, `a:link {  }`, `a:visited {  }`}
	for _, scss := range testCases {
		l := NewLexerWithString(scss)
		assert.NotNil(t, l)
		l.run()
		AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
		l.close()
	}
}

func TestLexerRuleLangPseudoSelector(t *testing.T) {
	// html:lang(fr-ca) { quotes: '« ' ' »' }
	l := NewLexerWithString(`html:lang(fr-ca) {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_FUNCTIONAL_PSEUDO, ast.T_PAREN_OPEN, ast.T_IDENT, ast.T_PAREN_CLOSE, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleLangPseudoSelectorAndFunctionalPseudoSelector(t *testing.T) {
	// html:lang(fr-ca) { quotes: '« ' ' »' }
	l := NewLexerWithString(`html:lang(fr-ca):first-child {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_TYPE_SELECTOR, ast.T_FUNCTIONAL_PSEUDO, ast.T_PAREN_OPEN, ast.T_IDENT, ast.T_PAREN_CLOSE,
		ast.T_PSEUDO_SELECTOR,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE,
	})
	l.close()
}

/*
func TestLexerRulePseudoSelectorWithNthChild(t *testing.T) {
	// html:lang(fr-ca) { quotes: '« ' ' »' }
	AssertLexerTokenSequence(t, `ui li:nth-child(3n+4) {  }`, []ast.TokenType{
		ast.T_TYPE_SELECTOR, ast.T_FUNCTIONAL_PSEUDO, ast.T_PAREN_OPEN, ast.T_IDENT, ast.T_PAREN_CLOSE,
		ast.T_PSEUDO_SELECTOR,
		ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE,
	})
}
*/

func TestLexerRuleWithIdSelector(t *testing.T) {
	AssertLexerTokenSequence(t, `#myPost {  }`, []ast.TokenType{ast.T_ID_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithTypeSelectorAndIdSelector(t *testing.T) {
	AssertLexerTokenSequence(t, `div#myPost {  }`, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_ID_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithTypeSelectorH1(t *testing.T) {
	AssertLexerTokenSequence(t, `h1 {  }`, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithTypeSelectorH1AndH2(t *testing.T) {
	AssertLexerTokenSequence(t, `h1, h2 {  }`, []ast.TokenType{ast.T_TYPE_SELECTOR, ast.T_COMMA, ast.T_TYPE_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerRuleWithIdSelectorWithDigits(t *testing.T) {
	AssertLexerTokenSequence(t, `#foo123 {  }`, []ast.TokenType{ast.T_ID_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerPropertyEmValueMul(t *testing.T) {
	l := NewLexerWithString(`.foo { width: 1.3em * 10.2em }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_FLOAT, ast.T_UNIT_EM, ast.T_MUL, ast.T_FLOAT, ast.T_UNIT_EM,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerPropertyPxValueMul(t *testing.T) {
	l := NewLexerWithString(`.foo { width: 1px * 3px }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_INTEGER, ast.T_UNIT_PX, ast.T_MUL, ast.T_INTEGER, ast.T_UNIT_PX,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithMultipleSelector(t *testing.T) {
	l := NewLexerWithString(`#foo123, .foo {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_ID_SELECTOR, ast.T_COMMA, ast.T_CLASS_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerRuleWithSubRuleWithParentSelector(t *testing.T) {
	l := NewLexerWithString(`.test {
		-webkit-transition: none;
		&.foo { color: #fff; }
	}`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_IDENT, ast.T_SEMICOLON,
		ast.T_PARENT_SELECTOR,
		ast.T_CLASS_SELECTOR,
		ast.T_BRACE_OPEN,
		ast.T_PROPERTY_NAME_TOKEN, ast.T_COLON, ast.T_HEX_COLOR, ast.T_SEMICOLON,
		ast.T_BRACE_CLOSE,
		ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorPseudoElementBefore(t *testing.T) {
	AssertLexerTokenSequence(t, `::before {  }`, []ast.TokenType{ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	AssertLexerTokenSequence(t, `::after {  }`, []ast.TokenType{ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	AssertLexerTokenSequence(t, `::first-line {  }`, []ast.TokenType{ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerSelectorInterpolationSuffix(t *testing.T) {
	AssertLexerTokenSequence(t, `#myPost#{ abc } {  }`, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerSelectorInterpolationPrefix(t *testing.T) {
	AssertLexerTokenSequence(t, `#{ abc }#myPost {  }`, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_ID_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
}

func TestLexerSelectorInterpolationWithPseudoSelector(t *testing.T) {
	l := NewLexerWithString(`#{ abc }:hover {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationWithPseudoSuffix(t *testing.T) {
	l := NewLexerWithString(`#{ abc }:hover {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_PSEUDO_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationInTheMiddleOfTypeSelector(t *testing.T) {
	l := NewLexerWithString(`foo#{ abc }bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationInTheMiddleOfClassSelector(t *testing.T) {
	l := NewLexerWithString(`.foo#{ abc }bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_CLASS_SELECTOR, ast.T_LITERAL_CONCAT, ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationWithSuffix(t *testing.T) {
	l := NewLexerWithString(`#{ abc }foo#{ bar } {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationInTheMiddleOfIdSelector(t *testing.T) {
	l := NewLexerWithString(`#foo#{ abc }bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationInTheMiddleOfPseudoSelector(t *testing.T) {
	l := NewLexerWithString(`:#{ abc }bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}

func TestLexerSelectorInterpolationInTheMiddleOfPseudoSelector2(t *testing.T) {
	l := NewLexerWithString(`:hover#{ abc }bar {  }`)
	assert.NotNil(t, l)
	l.run()
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_INTERPOLATION_SELECTOR, ast.T_BRACE_OPEN, ast.T_BRACE_CLOSE})
	l.close()
}
