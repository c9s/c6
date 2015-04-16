package c6

import "testing"
import "github.com/stretchr/testify/assert"

func TestLexer(t *testing.T) {
	l := NewLexerWithString(`.test {  }`)
	assert.NotNil(t, l)

	var r rune = l.next()
	assert.Equal(".", r)
}
