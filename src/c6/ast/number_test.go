package ast

import "testing"
import "github.com/stretchr/testify/assert"

func TestNumber(t *testing.T) {
	var num = NewIntegerNumber(100)
	num.AddInt(100)
	num.SetUnit(UNIT_PX)
	t.Logf("%s", num.String())

	assert.Equal(t, 200, num.Int)
	assert.Equal(t, "200px", num.String())
}
