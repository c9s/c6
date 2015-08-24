package ast

import "strings"

type ComplexSelectorList []*ComplexSelector

func (self *ComplexSelectorList) Append(sel *ComplexSelector) {
	var slice = append(*self, sel)
	*self = slice
}

func (self ComplexSelectorList) String() (css string) {
	var slices []string
	for _, sel := range self {
		slices = append(slices, sel.String())
	}
	return strings.Join(slices, ", ")
}
