package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestParser(t *testing.T) {
	parser := NewParser()

	assert.NotNil(t, parser)
	// assert.Equal(t, 123, 123, "they should be equal")
}
