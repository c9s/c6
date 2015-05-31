package ast

type RuleSet struct {
	Selectors *ComplexSelectorList
	Block     *DeclBlock
}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

func (self *RuleSet) AppendSubRuleSet(ruleset *RuleSet) {
	self.Block.AppendSubRuleSet(ruleset)
}

func (self *RuleSet) GetSubRuleSets() []*RuleSet {
	return self.Block.SubRuleSets
}

// Complete the statement interface
func (self *RuleSet) CanBeStmt() {}

func (self RuleSet) String() string {
	return "String() not implemented yet."
}
