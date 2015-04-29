package ast

/**
Variable can be used in block as statement and declaration block
*/
type VariableAssignment struct {
	Variable   *Variable
	Expression Expression
}

/*
Property is one of the declaration
*/
func (self VariableAssignment) CanBeDeclaration() {}
func (self VariableAssignment) CanBeStatement()   {}

func NewVariableAssignment(variable *Variable, expr Expression) *VariableAssignment {
	return &VariableAssignment{variable, expr}

}
