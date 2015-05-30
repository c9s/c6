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
	Stmts *StmtList

	// Nested rulesets
	SubRuleSets []*RuleSet

	Scope *Scope
}

func NewDeclarationBlock() *DeclarationBlock {
	return &DeclarationBlock{
		Stmts: &StmtList{},
	}
}

/**
Append a Declaration
*/
func (self *DeclarationBlock) Append(decl Stmt) {
	self.Stmts.Append(decl)
}

func (self *DeclarationBlock) AppendSubRuleSet(ruleset *RuleSet) {
	newRuleSets := append(self.SubRuleSets, ruleset)
	self.SubRuleSets = newRuleSets
}

func (self *DeclarationBlock) MergeStmts(stmts *StmtList) {
	for _, stm := range *stmts {
		self.Stmts.Append(stm)
	}
}

func (self DeclarationBlock) String() (out string) {
	out += "{\n"
	for _, decl := range *self.Stmts {
		out += decl.String() + "\n"
	}
	out += "}"
	return out
}
