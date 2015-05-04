package ast

import "testing"

// import "github.com/stretchr/testify/assert"

func TestBinaryExprString(t *testing.T) {
	var num1 = NewNumber(2, nil)
	var num2 = NewNumber(3, nil)
	var expr = NewBinaryExpression(OpAdd, num1, num2, false)
	t.Logf("%s", expr.String())

	var num3 = NewNumber(4, nil)
	var expr2 = NewBinaryExpression(OpSub, expr, num3, false)
	t.Logf("%s", expr2.String())
}
