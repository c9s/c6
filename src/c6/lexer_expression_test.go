package c6

import "testing"
import "c6/ast"

// import "github.com/stretchr/testify/assert"

func TestLexerExpression(t *testing.T) {
	l := NewLexerWithString(`$foo`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE})
	l.close()
}

func TestLexerHexColor3Letter(t *testing.T) {
	l := NewLexerWithString(`#fff`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_HEX_COLOR})
	l.close()
}

func TestLexerHexColor6Letter(t *testing.T) {
	l := NewLexerWithString(`#ffffff`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_HEX_COLOR})
	l.close()
}

func TestLexerFunctionCall(t *testing.T) {
	l := NewLexerWithString(`rgba(0,0,0,0)`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_FUNCTION_NAME, ast.T_PAREN_START,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER, ast.T_COMMA,
		ast.T_INTEGER,
		ast.T_PAREN_END})
	l.close()
}

func TestLexerExpressionFunction(t *testing.T) {
	l := NewLexerWithString(`rgba(0,0,0,0) - rgba(255,255,255,0)`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{
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
	l.close()
}

func TestLexerExpressionVariableMinusVariable(t *testing.T) {
	l := NewLexerWithString(`$foo - $bar`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_VARIABLE})
	l.close()
}

func TestLexerExpressionVariableNameWithDashSeparator(t *testing.T) {
	l := NewLexerWithString(`$a-b + 3px`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER, ast.T_UNIT_PX})
	l.close()
}

func TestLexerExpressionMinus3(t *testing.T) {
	l := NewLexerWithString(`$foo - 3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
	l.close()
}

func TestLexerExpressionMinus3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo-3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_MINUS, ast.T_INTEGER})
	l.close()
}

func TestLexerExpressionPlus3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo+3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_PLUS, ast.T_INTEGER})
	l.close()
}

func TestLexerExpressionDiv3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo/3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_DIV, ast.T_INTEGER})
	l.close()
}

func TestLexerExpressionMul3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo*3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []ast.TokenType{ast.T_VARIABLE, ast.T_MUL, ast.T_INTEGER})
	l.close()
}
