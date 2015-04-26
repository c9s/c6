package c6

import "github.com/stretchr/testify/assert"
import "testing"
import "fmt"
import "c6/ast"

func AssertLexerTokenSequence(t *testing.T, scss string, tokenList []ast.TokenType) {
	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	lexer.run()
	AssertTokenSequence(t, lexer, tokenList)
	lexer.close()
}

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []ast.TokenType) []ast.Token {
	fmt.Printf("Input: %s\n", l.Input)

	var tokens = []ast.Token{}
	var failure = false
	for _, expectingToken := range tokenList {

		var token = <-l.Output
		assert.NotNil(t, token)

		tokens = append(tokens, *token)

		if expectingToken == token.Type {
			fmt.Printf("\033[32mok %s '%s'\033[0m\n", token.Type.String(), token.Str)
		} else {
			failure = true
			fmt.Printf("\033[31mnot ok ---- got %s '%s' expecting %s\033[0m\n", token.Type.String(), token.Str, expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}
	if failure {
		t.FailNow()
	}

	return tokens
}

func AssertToken(t *testing.T, tokenType ast.TokenType, token *ast.Token) {
	assert.NotNil(t, token)
	if tokenType != token.Type {
		fmt.Printf("\033[31mnot ok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	} else {
		fmt.Printf("\033[32mok - expecting %s. Got %s '%s'\033[0m\n", tokenType.String(), token.Type.String(), token.Str)
	}
	assert.Equal(t, tokenType, token.Type)
}
