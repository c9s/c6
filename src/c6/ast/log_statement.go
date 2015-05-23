package ast

type LogStatement struct {
	Directive *Token
	Expr      Expression
}

func (stm LogStatement) CanBeStatement() {}
func (stm LogStatement) String() string  { return "LogStatement.String()" }
