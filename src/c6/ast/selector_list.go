package ast

/*
Gernal SelectorList struct
*/
type SelectorList []Selector

func (self *SelectorList) Append(sel Selector) {
	var slice = append(*self, sel)
	*self = slice
}

func (self SelectorList) Length() int {
	return len(self)
}

func (self SelectorList) String() (out string) {
	for _, sel := range self {
		out += sel.String()
	}
	return out
}

func NewSelectorList() *SelectorList {
	return &SelectorList{}
}
