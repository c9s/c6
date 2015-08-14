package lexer

import "github.com/stretchr/testify/assert"
import "testing"
import "github.com/c9s/c6/ast"

func AssertLexerTokenSequenceFromState(t *testing.T, scss string, fn stateFn, tokenList []ast.TokenType) {
	t.Logf("Testing SCSS: %s\n", scss)
	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	lexer.RunFrom(fn)
	AssertTokenSequence(t, lexer, tokenList)
	lexer.Close()
}

func AssertLexerTokenSequence(t *testing.T, scss string, tokenList []ast.TokenType) {
	t.Logf("Testing SCSS: %s\n", scss)
	var lexer = NewLexerWithString(scss)
	assert.NotNil(t, lexer)
	lexer.Run()
	AssertTokenSequence(t, lexer, tokenList)
	lexer.Close()
}

func OutputGreen(t *testing.T, msg string, args ...interface{}) {
	t.Logf("\033[32m")
	t.Logf(msg, args...)
	t.Logf("\033[0m\n")
}

func OutputRed(t *testing.T, msg string, args ...interface{}) {
	t.Logf("\033[31m")
	t.Logf(msg, args...)
	t.Logf("\033[0m\n")
}

func AssertTokenSequence(t *testing.T, l *Lexer, tokenList []ast.TokenType) []ast.Token {

	var tokens = []ast.Token{}
	var failure = false
	for idx, expectingToken := range tokenList {

		var token = <-l.Output

		if token == nil {
			failure = true
			OutputRed(t, "not ok ---- got nil expecting %s", expectingToken.String())
			break
		}

		tokens = append(tokens, *token)

		if expectingToken == token.Type {
			OutputGreen(t, "ok %s '%s'", token.Type.String(), token.Str)
		} else {
			failure = true
			OutputRed(t, "not ok ---- %d token => got %s '%s' expecting %s", idx, token.Type.String(), token.Str, expectingToken.String())
		}
		assert.Equal(t, expectingToken, token.Type)
	}

	if l.remaining() {
		var token *ast.Token = nil
		for token = <-l.Output; token != nil; token = <-l.Output {
			OutputRed(t, "not ok ---- Remaining expecting %s '%s'", token.Type.String(), token.Str)
		}
	}
	if failure {
		t.Fatal("See log.")
	}

	return tokens
}

func AssertTokenType(t *testing.T, tokenType ast.TokenType, token *ast.Token) {
	assert.NotNil(t, token)
	if tokenType != token.Type {
		OutputRed(t, "not ok - expecting %s. Got %s '%s'", tokenType.String(), token.Type.String(), token.Str)
	} else {
		OutputGreen(t, "ok - expecting %s. Got %s '%s'", tokenType.String(), token.Type.String(), token.Str)
	}
	assert.Equal(t, tokenType, token.Type)
}
