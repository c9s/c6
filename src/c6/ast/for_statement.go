package ast

type ForStatement struct {
	Variable *Variable
	From     Expression
	Through  Expression
	To       Expression
	Block    *Block
}

func (stm ForStatement) CanBeStatement() {}

func (stm ForStatement) String() string {
	return "@for " + stm.Variable.String() + " from " + stm.From.String() + " through " + stm.Through.String() + " {  }\n"
}

func NewForStatement(variable *Variable) *ForStatement {
	return &ForStatement{
		Variable: variable,
	}
}
