package ast

type Block struct {
	Scope *Scope
	Stmts *StmtList
}

type BlockNode interface {
	MergeStmts(stmts *StmtList)
}

func NewBlock() *Block {
	return &Block{
		Stmts: &StmtList{},
	}
}

// Override the statements
func (self *Block) SetStmts(stms *StmtList) {
	self.Stmts = stms
}

func (self *Block) MergeBlock(block *Block) {
	for _, stm := range *block.Stmts {
		*self.Stmts = append(*self.Stmts, stm)
	}
}

func (self *Block) MergeStmts(stmts *StmtList) {
	for _, stm := range *stmts {
		self.Stmts.Append(stm)
	}
}
