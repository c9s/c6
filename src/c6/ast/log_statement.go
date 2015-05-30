package ast

type LogStmt struct {
	Directive *Token
	Expr      Expression
}

func (stm LogStmt) CanBeStmt() {}
func (stm LogStmt) String() string  { return "LogStmt.String()" }
