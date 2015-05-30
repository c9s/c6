package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestBinaryExprString(t *testing.T) {
	var num1 = NewNumber(2, NewUnit(T_UNIT_PX, nil), nil)
	var num2 = NewNumber(3, NewUnit(T_UNIT_PX, nil), nil)
	var expr = NewBinaryExpr(NewOp(T_DIV), num1, num2, false)
	assert.Equal(t, "2px/3px", expr.String())
}
