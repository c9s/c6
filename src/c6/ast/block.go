package ast

import "github.com/c9s/c6/src/c6/symtable"

type Block struct {
	SymTable   *symtable.SymTable
	Statements *StatementList
}

type BlockNode interface {
	MergeStatements(stmts *StatementList)
	GetSymTable() *symtable.SymTable
}

func NewBlock() *Block {
	return &Block{
		SymTable:   &symtable.SymTable{},
		Statements: &StatementList{},
	}
}

func (self *Block) GetSymTable() *symtable.SymTable {
	return self.SymTable
}

// Override the statements
func (self *Block) SetStatements(stms *StatementList) {
	self.Statements = stms
}

func (self *Block) MergeBlock(block *Block) {
	for _, stm := range *block.Statements {
		*self.Statements = append(*self.Statements, stm)
	}
}

func (self *Block) MergeStatements(stmts *StatementList) {
	for _, stm := range *stmts {
		self.Statements.Append(stm)
	}
}
