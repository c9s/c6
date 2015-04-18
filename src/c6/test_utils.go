package c6

import "github.com/stretchr/testify/assert"
import "testing"

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []TokenType) {
	for _, expectingToken := range tokenList {
		var token = <-l.Output
		assert.NotNil(t, token)

		if expectingToken == token.Type {
			t.Logf("Expecting Token: %s - ok", expectingToken.String())
		} else {
			t.Logf("Expecting Token: %s - not ok", expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}
}
