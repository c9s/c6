package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestParserParseEmptyRuleWithClassSelector(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
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
