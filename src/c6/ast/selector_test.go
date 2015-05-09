package ast

import "github.com/stretchr/testify/assert"
import "testing"

func TestCombinedSelector(t *testing.T) {
	e := TypeSelector{"div"}
	cls1 := ClassSelector{".foo"}
	cls2 := ClassSelector{".bar"}
	id := IdSelector{"#myId"}

	assert.Equal(t, ".foo", cls1.String())
	assert.Equal(t, ".bar", cls2.String())
	assert.Equal(t, "#myId", id.String())

	combined := CombinedSelector{"", []Selector{e, id, cls1, cls2}}
	assert.Equal(t, "div#myId.foo.bar", combined.String())
}
