package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestParserParseEmptyRuleWithClassSelector(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
	parser.parseScss(`.test {  }`)

	var token *Token

	token = parser.next()
	assert.NotNil(t, token)
	assert.Equal(t, T_CLASS_SELECTOR, token.Type)

	token = parser.next()
	assert.NotNil(t, token)
	assert.Equal(t, T_BRACE_START, token.Type)

	token = parser.next()
	assert.NotNil(t, token)
	assert.Equal(t, T_BRACE_END, token.Type)
	// assert.Equal(t, 123, 123, "they should be equal")
}
