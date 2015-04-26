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

func (self *RuleSet) AppendDeclaration(dec Declaration) {
}

func (self *RuleSet) AppendSubRuleSet() {

}

// Complete the statement interface
func (self *RuleSet) CanBeStatement() {}
