package ast

/**
Variable can be used in block as statement and declaration block
*/
type VariableAssignment struct {
	Variable   Variable
	Expression Expression
}

/*
Property is one of the declaration
*/
func (self VariableAssignment) IsDeclaration() {}
func (self VariableAssignment) IsStatement()   {}
