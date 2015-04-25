package ast

type Block struct {
	Statements []Statement
}

// Override the statements
func (self *Block) SetStatements(stms []Statement) {
	self.Statements = stms
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
