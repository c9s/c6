package c6

import "github.com/stretchr/testify/assert"
import "testing"
import "fmt"

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []TokenType) {
	fmt.Printf("Input: %s\n", l.Input)
	for _, expectingToken := range tokenList {

		var token = <-l.Output
		assert.NotNil(t, token)

		if expectingToken == token.Type {
			fmt.Printf("\033[32mok %s '%s'\033[0m\n", token.Type.String(), token.Str)
		} else {
			fmt.Printf("\033[31mnot ok ---- got %s '%s' expecting %s\033[0m\n", token.Type.String(), token.Str, expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}
}

func AssertToken(t *testing.T, tokenType TokenType, token *Token) {
	assert.NotNil(t, token)
	if tokenType != token.Type {
		fmt.Printf("\033[31mnot ok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	} else {
		fmt.Printf("\033[32mok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	}
	assert.Equal(t, tokenType, token.Type)
}
