package ast

type RuleSet struct {
	Selectors        []Selector
	DeclarationBlock *DeclarationBlock
}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

func (self *RuleSet) AppendSelector(sel Selector) {
	self.Selectors = append(self.Selectors, sel)
}

func (self *RuleSet) AppendSubRuleSet(ruleset *RuleSet) {
	self.DeclarationBlock.AppendSubRuleSet(ruleset)
}

func (self *RuleSet) GetSubRuleSets() []*RuleSet {
	return self.DeclarationBlock.SubRuleSets
}

// Complete the statement interface
func (self *RuleSet) CanBeStatement() {}

func (self RuleSet) String() string {
	return "String() not implemented yet."
}
