package ast

import "github.com/c9s/c6/src/c6/symtable"

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	SymTable *symtable.SymTable

	Statements *StatementList

	// Nested rulesets
	SubRuleSets []*RuleSet
}

func NewDeclarationBlock() *DeclarationBlock {
	return &DeclarationBlock{
		Statements: &StatementList{},
		SymTable:   &symtable.SymTable{},
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

func (self *DeclarationBlock) GetSymTable() *symtable.SymTable {
	return self.SymTable
}

func (self DeclarationBlock) String() (out string) {
	out += "{\n"
	for _, decl := range *self.Statements {
		out += decl.String() + "\n"
	}
	out += "}"
	return out
}
