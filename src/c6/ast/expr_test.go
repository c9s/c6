package ast

import "testing"

// import "github.com/stretchr/testify/assert"

func TestBinaryExprString(t *testing.T) {
	var num1 = NewNumber(2, NewUnit(T_UNIT_PX, nil), nil)
	var num2 = NewNumber(3, NewUnit(T_UNIT_PX, nil), nil)
	var expr = NewBinaryExpression(NewOp(T_DIV, nil), num1, num2, false)
	t.Logf("%s", expr.String())
}
