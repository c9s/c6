package ast

type ComplexSelectorItem struct {
	Combinator       Combinator
	CompoundSelector *CompoundSelector
}

func (item *ComplexSelectorItem) String() (css string) {
	return item.Combinator.String() + item.CompoundSelector.String()
}
func (item *ComplexSelectorItem) CSSString() (css string) { return item.String() }

type ComplexSelector struct {
	CompoundSelector     *CompoundSelector
	ComplexSelectorItems []*ComplexSelectorItem
}

func (self *ComplexSelector) AppendCompoundSelector(comb Combinator, sel *CompoundSelector) {
	self.ComplexSelectorItems = append(self.ComplexSelectorItems, &ComplexSelectorItem{comb, sel})
}

func (self *ComplexSelector) String() (css string) {
	css = self.CompoundSelector.String()
	for _, item := range self.ComplexSelectorItems {
		css += item.Combinator.String()
		css += item.CompoundSelector.String()
	}
	return css
}

func (self *ComplexSelector) CSSString() string { return self.String() }

func NewComplexSelector(sel *CompoundSelector) *ComplexSelector {
	return &ComplexSelector{
		CompoundSelector: sel,
	}
}
