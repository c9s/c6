package ast

import "c6/symtable"

type Block struct {
	SymTable   *symtable.SymTable
	Statements *StatementList
}

func NewBlock() *Block {
	return &Block{}
}

// Override the statements
func (self *Block) SetStatements(stms *StatementList) {
	self.Statements = stms
}

func (self *Block) MergeBlock(block *Block) {
	for _, stm := range *block.Statements {
		self.Statements.Append(stm)
	}
}

func (self *Block) AppendStatements(stmts []Statement) {
	for _, stm := range stmts {
		self.Statements.Append(stm)
	}
}
