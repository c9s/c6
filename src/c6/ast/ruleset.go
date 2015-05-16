package ast

type SelectorList []Selector

func (self SelectorList) Append(sel Selector) {
	newSlice := append(self, sel)
	self = newSlice
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

type RuleSet struct {
	Selectors *SelectorList
	Block     *DeclarationBlock
}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

func (self *RuleSet) AppendSelector(sel Selector) {
	self.Selectors.Append(sel)
}

func (self *RuleSet) AppendSubRuleSet(ruleset *RuleSet) {
	self.Block.AppendSubRuleSet(ruleset)
}

func (self *RuleSet) GetSubRuleSets() []*RuleSet {
	return self.Block.SubRuleSets
}

// Complete the statement interface
func (self *RuleSet) CanBeStatement() {}

func (self RuleSet) String() string {
	return "String() not implemented yet."
}
