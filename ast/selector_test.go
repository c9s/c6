package ast

import "github.com/stretchr/testify/assert"
import "testing"

func TestSelectorList(t *testing.T) {

	e := NewTypeSelector("div")
	cls1 := NewClassSelector(".foo")
	cls2 := NewClassSelector(".bar")
	id := NewIdSelector("#myId")

	assert.Equal(t, ".foo", cls1.String())
	assert.Equal(t, ".bar", cls2.String())
	assert.Equal(t, "#myId", id.String())

	combined := SelectorList{e, id, cls1, cls2}
	assert.Equal(t, "div#myId.foo.bar", combined.String())
}
