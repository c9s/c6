package ast

type Block struct {
	Statements *StatementList
}

type BlockNode interface {
	MergeStatements(stmts *StatementList)
}

func NewBlock() *Block {
	return &Block{
		Statements: &StatementList{},
	}
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
