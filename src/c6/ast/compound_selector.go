package ast

// one or more simple selector
type CompoundSelector []Selector

func (self *CompoundSelector) Append(sel Selector) {
	var slice = append(*self, sel)
	*self = slice
}

func (self CompoundSelector) Length() int {
	return len(self)
}

func (self CompoundSelector) String() (css string) {
	for _, sel := range self {
		css += sel.String()
	}
	return css
}

func NewCompoundSelector() *CompoundSelector {
	return &CompoundSelector{}
}
