package ast

import _ "github.com/stretchr/testify/assert"
import "testing"

func TestCombinedSelector(t *testing.T) {
	e := TypeSelector{"div"}
	cls1 := ClassSelector{"foo"}
	cls2 := ClassSelector{"bar"}
	combined := CombinedSelector{"", []Selector{Selector(e), Selector(cls1), Selector(cls2)}}
	_ = combined
}
