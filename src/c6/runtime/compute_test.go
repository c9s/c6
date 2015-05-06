package runtime

import "c6/ast"
import "testing"
import "github.com/stretchr/testify/assert"

func TestComputeNumberAddNumber(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_PLUS, nil), ast.NewNumber(10, nil, nil), ast.NewNumber(3, nil, nil))
	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, 13.0, num.Value)
}

func TestComputeNumberAddNumberIncompatibleUnit(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_PLUS, nil), ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil), ast.NewNumber(3, ast.NewUnit(ast.T_UNIT_PT, nil), nil))
	assert.Nil(t, val)
}

func TestComputeNumberMulWithUnit(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_MUL, nil), ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil), ast.NewNumber(3, nil, nil))
	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, ast.T_UNIT_PX, num.Unit.Type)
	assert.Equal(t, 30.0, num.Value)
}

func TestComputeNumberDivWithUnit(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_DIV, nil),
		ast.NewNumber(10, ast.NewUnit(ast.T_UNIT_PX, nil), nil),
		ast.NewNumber(2, nil, nil))

	num, ok := val.(*ast.Number)
	assert.True(t, ok)
	assert.NotNil(t, num.Unit)
	assert.Equal(t, ast.T_UNIT_PX, num.Unit.Type)
	assert.Equal(t, 5.0, num.Value)
}

func TestComputeRGBAColorWithNumber(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_PLUS, nil), ast.NewRGBAColor(10, 10, 10, 0.2, nil), ast.NewNumber(3, nil, nil))
	c, ok := val.(*ast.RGBAColor)
	assert.True(t, ok)
	assert.Equal(t, "rgba(13, 13, 13, 0.2)", c.String())
}

func TestComputeRGBColorWithNumber(t *testing.T) {
	val := Compute(ast.NewOp(ast.T_PLUS, nil), ast.NewRGBColor(10, 10, 10, nil), ast.NewNumber(3, nil, nil))
	c, ok := val.(*ast.RGBColor)
	assert.True(t, ok)
	assert.Equal(t, "rgb(13, 13, 13)", c.String())
}
