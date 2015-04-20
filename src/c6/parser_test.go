package c6

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func AssertToken(t *testing.T, tokenType TokenType, token *Token) {
	assert.NotNil(t, token)
	if tokenType != token.Type {
		fmt.Printf("\033[31mnot ok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	} else {
		fmt.Printf("\033[32mok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	}
	assert.Equal(t, tokenType, token.Type)
}

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
