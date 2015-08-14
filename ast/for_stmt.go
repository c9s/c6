package ast

type ForStmt struct {
	Variable *Variable
	From     Expr
	Through  Expr
	To       Expr
	Block    *DeclBlock
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
