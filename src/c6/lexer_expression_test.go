package c6

import "testing"

// import "github.com/stretchr/testify/assert"

func TestLexerExpression(t *testing.T) {
	l := NewLexerWithString(`$foo`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE})
	l.close()
}

func TestLexerExpressionVariableMinusVariable(t *testing.T) {
	l := NewLexerWithString(`$foo - $bar`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_MINUS, T_VARIABLE})
	l.close()
}

func TestLexerExpressionVariableNameWithDashSeparator(t *testing.T) {
	l := NewLexerWithString(`$a-b + 3px`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_PLUS, T_INTEGER, T_UNIT_PX})
	l.close()
}

func TestLexerExpressionMinus3(t *testing.T) {
	l := NewLexerWithString(`$foo - 3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_MINUS, T_INTEGER})
	l.close()
}

func TestLexerExpressionMinus3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo-3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_MINUS, T_INTEGER})
	l.close()
}

func TestLexerExpressionPlus3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo+3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_PLUS, T_INTEGER})
	l.close()
}

func TestLexerExpressionDiv3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo/3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_DIV, T_INTEGER})
	l.close()
}

func TestLexerExpressionMul3WithoutSpace(t *testing.T) {
	l := NewLexerWithString(`$foo*3`)
	l.runFrom(lexExpression)
	AssertTokenSequence(t, l, []TokenType{T_VARIABLE, T_MUL, T_INTEGER})
	l.close()
}
