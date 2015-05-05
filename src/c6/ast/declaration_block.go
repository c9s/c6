package ast

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	Declarations []Declaration

	// Nested rulesets
	SubRuleSets []*RuleSet
}

/**
Append a Declaration
*/
func (self *DeclarationBlock) Append(decl Declaration) {
	self.Declarations = append(self.Declarations, decl)
}

func (self *DeclarationBlock) AppendSubRuleSet(ruleset *RuleSet) {
	newRuleSets := append(self.SubRuleSets, ruleset)
	self.SubRuleSets = newRuleSets
}

func (self DeclarationBlock) String() (out string) {
	out += "{\n"
	for _, decl := range self.Declarations {
		out += decl.String() + "\n"
	}
	out += "}"
	return out
}
