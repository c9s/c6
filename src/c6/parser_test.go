package c6

import "testing"
import "c6/ast"
import "github.com/stretchr/testify/assert"

import "fmt"

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
	assert.Equal(t, "http://foo.com/bar.css", string(url))
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
	assert.Equal(t, "foo.css", string(url))
}

func TestParserParseImportRuleWithMediaList(t *testing.T) {
	parser := NewParser()
	block := parser.parseScss(`@import url("foo.css") screen;`)
	_ = block
}

func TestParserPropertyListExpression(t *testing.T) {
	var buffers []string = []string{
		`div { width: 1px; }`,
		`div { width: 2px 3px; }`,
		`div { width: 4px, 5px, 6px, 7px; }`,
		`div { width: 4px, 5px 6px, 7px; }`,
		`div { width: 10px 3px + 7px 20px; }`,
		// `div { width: 10px, 3px + 7px, 20px; }`,
	}
	for _, buffer := range buffers {
		t.Logf("Input %s", buffer)
		var parser = NewParser()
		var block = parser.parseScss(buffer)
		fmt.Printf("%+v\n", block)
	}
}

func TestParserVariableAssignmentWithMorePlus(t *testing.T) {
	var parser = NewParser()
	var block = parser.parseScss(`$foo: 12px + 20px + 20px;`)
	fmt.Printf("%+v\n", block)
}

func TestParserVariableAssignmentWithSimpleExpresion(t *testing.T) {
	var parser = NewParser()
	var block = parser.parseScss(`$foo: 10px + 20px;`)
	fmt.Printf("%+v\n", block)
}

func TestParserVariableAssignmentWithPxValue(t *testing.T) {
	var parser = NewParser()
	var block = parser.parseScss(`$foo: 10px;`)
	fmt.Printf("%+v\n", block)
}

func TestParserMassiveRules(t *testing.T) {
	var buffers []string = []string{
		`div { width: auto; }`,
		`div { width: 100px }`,
		`div { width: 100pt }`,
		`div { width: 100em }`,
		`div { width: 100rem }`,
		`div { padding: 10px 10px; }`,
		`div { padding: 10px 10px 20px 30px; }`,
		`div { padding: 10px + 10px; }`,
		`div { padding: 10px + 10px * 3; }`,
		`div { color: red; }`,
		`div { color: rgb(255,255,255); }`,
		`div { color: rgba(255,255,255,0); }`,
		`div { background-image: url("../images/foo.png"); }`,
		// `div { color: #ccddee; }`,
	}
	for _, buffer := range buffers {
		t.Logf("Input %s", buffer)
		var parser = NewParser()
		var block = parser.parseScss(buffer)
		fmt.Printf("%+v\n", block)
	}
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
	AssertTokenType(t, T_CLASS_SELECTOR, token)

	// should be the same
	token = parser.next()
	AssertTokenType(t, T_CLASS_SELECTOR, token)

	token = parser.next()
	AssertTokenType(t, T_BRACE_START, token)

	token = parser.peek()
	AssertTokenType(t, T_BRACE_END, token)

	token = parser.next()
	AssertTokenType(t, T_BRACE_END, token)
}

func TestParserPeekBy(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
	parser.parseScss(`.test {  }`)

	var token *Token

	token = parser.peekBy(0)
	AssertTokenType(t, T_CLASS_SELECTOR, token)

	token = parser.peekBy(1)
	AssertTokenType(t, T_BRACE_START, token)

	token = parser.peekBy(2)
	AssertTokenType(t, T_BRACE_END, token)

	token = parser.next()
	AssertTokenType(t, T_CLASS_SELECTOR, token)
	token = parser.next()
	AssertTokenType(t, T_BRACE_START, token)
	token = parser.next()
	AssertTokenType(t, T_BRACE_END, token)
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
