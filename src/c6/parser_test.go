package c6

import "testing"
import "c6/ast"
import "github.com/stretchr/testify/assert"

func TestParserParseImportRuleWithUrl(t *testing.T) {
	parser := NewParser()
	block := parser.parseScss(`@import url("http://foo.com/bar.css");`)

	rule, ok := block.Statement(0).(*ast.ImportStatement)
	assert.True(t, ok, "Convert to ImportStatement OK")
	assert.NotNil(t, rule)

	// it's not a relative url
	_, ok1 := rule.Url.(ast.RelativeUrl)
	assert.False(t, ok1)

	// it's a url
	url, ok2 := rule.Url.(ast.Url)
	assert.True(t, ok2)
	assert.Equal(t, "http://foo.com/bar.css", url)
}

func TestParserParseImportRuleWithString(t *testing.T) {
	parser := NewParser()
	block := parser.parseScss(`@import "foo.css";`)

	rule, ok := block.Statement(0).(*ast.ImportStatement)
	assert.True(t, ok, "Convert to ImportStatement OK")
	assert.NotNil(t, rule)

	// it's not a relative url
	url, ok := rule.Url.(ast.RelativeUrl)
	assert.True(t, ok)

	assert.True(t, ok)
	assert.Equal(t, "foo.css", url)
}

func TestParserParseImportRuleWithMediaList(t *testing.T) {
	parser := NewParser()
	block := parser.parseScss(`@import url("foo.css") screen;`)
	_ = block
}

func TestParserParseTypeSelectorRule(t *testing.T) {
	parser := NewParser()
	block := parser.parseScss(`div { width: auto; }`)

	ruleset, ok := block.Statements[0].(*ast.RuleSet)
	assert.True(t, ok)

	t.Logf("%+v\n", ruleset.Selectors)
	t.Logf("%+v\n", ruleset.DeclarationBlock)

	// _ = block
}

/*
func TestParserParseEmptyRuleWithClassSelector(t *testing.T) {
	parser := NewParser()
	parser.parseScss(`.test {  }`)

	var token *Token

	token = parser.peek()
	AssertToken(t, T_CLASS_SELECTOR, token)

	// should be the same
	token = parser.next()
	AssertToken(t, T_CLASS_SELECTOR, token)

	token = parser.next()
	AssertToken(t, T_BRACE_START, token)

	token = parser.peek()
	AssertToken(t, T_BRACE_END, token)

	token = parser.next()
	AssertToken(t, T_BRACE_END, token)
}

func TestParserPeekBy(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
	parser.parseScss(`.test {  }`)

	var token *Token

	token = parser.peekBy(0)
	AssertToken(t, T_CLASS_SELECTOR, token)

	token = parser.peekBy(1)
	AssertToken(t, T_BRACE_START, token)

	token = parser.peekBy(2)
	AssertToken(t, T_BRACE_END, token)

	token = parser.next()
	AssertToken(t, T_CLASS_SELECTOR, token)
	token = parser.next()
	AssertToken(t, T_BRACE_START, token)
	token = parser.next()
	AssertToken(t, T_BRACE_END, token)
}

func TestParseNestedRule(t *testing.T) {
	code := `
#main p {
  color: #00ff00;
  width: 97%;

  .redbox {
    background-color: #ff0000;
    color: #000000;
  }
}
`
	p := NewParser()
	assert.NotNil(t, p)
	p.parseScss(code)
}
*/
