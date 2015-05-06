package ast

type Block struct {
	Statements []Statement
}

func NewBlock() *Block {
	return &Block{}
}

// Override the statements
func (self *Block) SetStatements(stms []Statement) {
	self.Statements = stms
}

func (self *Block) MergeBlock(block *Block) {
	for _, stm := range block.Statements {
		self.Statements = append(self.Statements, stm)
	}
}

// Append a statement interface object on the slice
// It needs to be a pointer
func (self *Block) AppendStatement(stm Statement) {
	self.Statements = append(self.Statements, stm)
}

// Return the reference of the statement
func (self *Block) Statement(idx uint) Statement {
	return self.Statements[idx]
}

/*
func (self Block) String() string {
	for _, statement := range self.Statements {
		switch t := statement.(type) {
			case *VariabeAssignment:
				t.String()
		}
	}
}
*/
