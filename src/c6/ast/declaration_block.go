package ast

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	// The symbol table for storing constant values
	// Only constants can be stored here...
	Statements *StatementList

	// Nested rulesets
	SubRuleSets []*RuleSet
}

func NewDeclarationBlock() *DeclarationBlock {
	return &DeclarationBlock{
		Statements: &StatementList{},
	}
}

/**
Append a Declaration
*/
func (self *DeclarationBlock) Append(decl Statement) {
	self.Statements.Append(decl)
}

func (self *DeclarationBlock) AppendSubRuleSet(ruleset *RuleSet) {
	newRuleSets := append(self.SubRuleSets, ruleset)
	self.SubRuleSets = newRuleSets
}

func (self *DeclarationBlock) MergeStatements(stmts *StatementList) {
	for _, stm := range *stmts {
		self.Statements.Append(stm)
	}
}

func (self DeclarationBlock) String() (out string) {
	out += "{\n"
	for _, decl := range *self.Statements {
		out += decl.String() + "\n"
	}
	out += "}"
	return out
}
