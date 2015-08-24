package ast

/**
Variable can be used in block as statement and declaration block
*/
type AssignStmt struct {
	Variable *Variable
	Expr     Expr

	Default   bool
	Optional  bool
	Global    bool
	Important bool
}

/*
Property is one of the declaration
*/
func (self AssignStmt) CanBeDeclaration() {}
func (self AssignStmt) CanBeStmt()        {}

func (self AssignStmt) String() string {
	return self.Variable.String() + " = " + self.Expr.String()
}

func NewAssignStmt(variable *Variable, expr Expr) *AssignStmt {
	return &AssignStmt{variable, expr, false, false, false, false}
}
