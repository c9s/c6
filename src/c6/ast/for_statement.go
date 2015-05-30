package ast

type ForStmt struct {
	Variable *Variable
	From     Expression
	Through  Expression
	To       Expression
	Block    *DeclarationBlock
}

func (stm ForStmt) CanBeStmt() {}

func (stm ForStmt) String() string {
	return "@for " + stm.Variable.String() + " from " + stm.From.String() + " through " + stm.Through.String() + " {  }\n"
}

func NewForStmt(variable *Variable) *ForStmt {
	return &ForStmt{
		Variable: variable,
	}
}
