package c6

import "testing"
import "c6/ast"

// import "github.com/stretchr/testify/assert"

func TestLexerIdentifier(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `none`, lexExpression2, []ast.TokenType{ast.T_IDENT})
}

func TestLexerIdentifierWithTrailingInterp(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `none#{ 10 + 10 }`, lexExpression2, []ast.TokenType{
		ast.T_IDENT, ast.T_LITERAL_CONCAT, ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END})
}

func TestLexerIdentifierWithLeadingInterp(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#{ 10 + 10 }none`, lexExpression2, []ast.TokenType{
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END, ast.T_LITERAL_CONCAT, ast.T_IDENT})
}

func TestLexerExpression(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo`, lexExpression2, []ast.TokenType{ast.T_VARIABLE})
}

func TestLexerHexColor3Letter(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#fff`, lexExpression2, []ast.TokenType{ast.T_HEX_COLOR})
}

func TestLexerHexColor6Letter(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#ffffff`, lexExpression2, []ast.TokenType{ast.T_HEX_COLOR})
}

func TestLexerFunctionCall(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `rgba(0,0,0,0)`, lexExpression2, []ast.TokenType{ast.T_FUNCTION_NAME, ast.T_PAREN_START,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER,
		ast.T_PAREN_END})
}

func TestLexerExpressionFunction(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `rgba(0,0,0,0) - rgba(255,255,255,0)`, lexExpression2, []ast.TokenType{
		ast.T_FUNCTION_NAME, ast.T_PAREN_START,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_PAREN_END,
		ast.T_MINUS,
		ast.T_FUNCTION_NAME, ast.T_PAREN_START,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_PAREN_END,
	})
}

func TestLexerExpressionUnicodeRange1(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+0025-00FF`, lexExpression2, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExpressionUnicodeRange2(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+26`, lexExpression2, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExpressionUnicodeRange3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+0025`, lexExpression2, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExpressionVariableMinusVariable(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo - $bar`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_VARIABLE})
}

func TestLexerExpressionVariableNameWithDashSeparator(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$a-b + 3px`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER, ast.T_UNIT_PX})
}

func TestLexerExpressionMinus3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo - 3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExpressionMinus43(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-4-3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExpressionMinus4neg3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-4--3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER, ast.T_MINUS, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExpressionMinus3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExpressionPlus3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo+3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER})
}

func TestLexerExpressionDiv3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo/3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_DIV, ast.T_INTEGER})
}

func TestLexerExpressionMul3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo*3`, lexExpression2, []ast.TokenType{ast.T_VARIABLE, ast.T_MUL, ast.T_INTEGER})
}
