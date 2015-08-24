package runtime

import "testing"
import "github.com/c9s/c6/ast"
import "github.com/stretchr/testify/assert"

func TestStack(t *testing.T) {
	var stack = &Stack{}
	stack.Push(ast.NewNumber(10, nil, nil))
	stack.Push(ast.NewNumber(12, nil, nil))

	val1 := stack.Pop()
	num1, ok := val1.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, 12.0, num1.Value)

	val2 := stack.Pop()
	num2, ok := val2.(*ast.Number)
	assert.True(t, ok)
	assert.Equal(t, 10.0, num2.Value)
}
