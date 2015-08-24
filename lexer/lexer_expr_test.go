package lexer

import "testing"
import "github.com/c9s/c6/ast"

// import "github.com/stretchr/testify/assert"

func TestLexerIdentifier(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `none`, lexExpr, []ast.TokenType{ast.T_IDENT})
}

func TestLexerIdentifierWithTrailingInterp(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `none#{ 10 + 10 }`, lexExpr, []ast.TokenType{
		ast.T_IDENT, ast.T_LITERAL_CONCAT, ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END})
}

func TestLexerIdentifierWithLeadingInterp(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#{ 10 + 10 }none`, lexExpr, []ast.TokenType{
		ast.T_INTERPOLATION_START,
		ast.T_INTEGER, ast.T_PLUS, ast.T_INTEGER,
		ast.T_INTERPOLATION_END, ast.T_LITERAL_CONCAT, ast.T_IDENT})
}

func TestLexerExpr(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo`, lexExpr, []ast.TokenType{ast.T_VARIABLE})
}

func TestLexerHexColor3Letter(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#fff`, lexExpr, []ast.TokenType{ast.T_HEX_COLOR})
}

func TestLexerHexColor6Letter(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `#ffffff`, lexExpr, []ast.TokenType{ast.T_HEX_COLOR})
}

func TestLexerFunctionCall(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `rgba(0,0,0,0)`, lexExpr, []ast.TokenType{ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER,
		ast.T_PAREN_CLOSE})
}

func TestLexerExprFunction(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `rgba(0,0,0,0) - rgba(255,255,255,0)`, lexExpr, []ast.TokenType{
		ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_PAREN_CLOSE,
		ast.T_MINUS,
		ast.T_FUNCTION_NAME, ast.T_PAREN_OPEN,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_PAREN_CLOSE,
	})
}

func TestLexerExprUnicodeRange1(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+0025-00FF`, lexExpr, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExprUnicodeRange2(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+26`, lexExpr, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExprUnicodeRange3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `U+0025`, lexExpr, []ast.TokenType{ast.T_UNICODE_RANGE})
}

func TestLexerExprNumberWithScientificNotation(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `10e3`, lexExpr, []ast.TokenType{ast.T_INTEGER})
}

func TestLexerExprNumberWithScientificNotation2(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `3.4e-2`, lexExpr, []ast.TokenType{ast.T_FLOAT})
}

func TestLexerExprNumberWithNthChildSyntax(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `3n+1`, lexExpr, []ast.TokenType{ast.T_INTEGER, ast.T_N, ast.T_PLUS, ast.T_INTEGER})
}

func TestLexerExprNumberWithNthChildSyntax2(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `3n`, lexExpr, []ast.TokenType{ast.T_INTEGER, ast.T_N})
}

func TestLexerExprNumberWithNthChildSyntax3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `-n`, lexExpr, []ast.TokenType{ast.T_MINUS, ast.T_N})
}

func TestLexerExprNumberWithNthChildSyntax4(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `2n`, lexExpr, []ast.TokenType{ast.T_INTEGER, ast.T_N})
}

func TestLexerExprNumberWithScientificNotation3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `3.1415e-5`, lexExpr, []ast.TokenType{ast.T_FLOAT})
}

func TestLexerExprVariableMinusVariable(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo - $bar`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_VARIABLE})
}

func TestLexerExprVariableNameWithDashSeparator(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$a-b + 3px`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER, ast.T_UNIT_PX})
}

func TestLexerExprMinus3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo - 3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExprMinus43(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-4-3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExprMinus4neg3(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-4--3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER, ast.T_MINUS, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExprMinus3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo-3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
}

func TestLexerExprPlus3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo+3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER})
}

func TestLexerExprDiv3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo/3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_DIV, ast.T_INTEGER})
}

func TestLexerExprMul3WithoutSpace(t *testing.T) {
	AssertLexerTokenSequenceFromState(t, `$foo*3`, lexExpr, []ast.TokenType{ast.T_VARIABLE, ast.T_MUL, ast.T_INTEGER})
}
