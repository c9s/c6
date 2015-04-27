package c6

import "github.com/stretchr/testify/assert"
import "testing"
import "fmt"
import "c6/ast"

func AssertLexerTokenSequence(t *testing.T, scss string, tokenList []ast.TokenType) {
	fmt.Printf("Testing SCSS: %s\n", scss)

	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	lexer.run()
	AssertTokenSequence(t, lexer, tokenList)
	lexer.close()
}

func OutputGreen(msg string, args ...interface{}) {
	fmt.Printf("\033[32m")
	fmt.Printf(msg, args...)
	fmt.Printf("\033[0m\n")
}

func OutputRed(msg string, args ...interface{}) {
	fmt.Printf("\033[31m")
	fmt.Printf(msg, args...)
	fmt.Printf("\033[0m\n")
}

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []ast.TokenType) []ast.Token {

	var tokens = []ast.Token{}
	var failure = false
	for _, expectingToken := range tokenList {

		var token = <-l.Output
		assert.NotNil(t, token)

		tokens = append(tokens, *token)

		if expectingToken == token.Type {
			OutputGreen("ok %s '%s'", token.Type.String(), token.Str)
		} else {
			failure = true
			OutputRed("not ok ---- got %s '%s' expecting %s", token.Type.String(), token.Str, expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}

	if l.remaining() {
		var token *ast.Token = nil
		for token = <-l.Output; token != nil; token = <-l.Output {
			OutputRed("not ok ---- Remaining expecting %s '%s'", token.Type.String(), token.Str)
		}
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
