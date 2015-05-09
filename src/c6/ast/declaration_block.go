package ast

import "c6/symtable"

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	SymTable *symtable.SymTable

	Statements []Statement

	// Nested rulesets
	SubRuleSets []*RuleSet
}

/**
Append a Declaration
*/
func (self *DeclarationBlock) Append(decl Statement) {
	self.Statements = append(self.Statements, decl)
}

func (self *DeclarationBlock) AppendSubRuleSet(ruleset *RuleSet) {
	newRuleSets := append(self.SubRuleSets, ruleset)
	self.SubRuleSets = newRuleSets
}

func (self DeclarationBlock) String() (out string) {
	out += "{\n"
	for _, decl := range self.Statements {
		out += decl.String() + "\n"
	}
	out += "}"
	return out
}
