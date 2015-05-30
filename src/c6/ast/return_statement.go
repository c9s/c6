package ast

type ReturnStmt struct {
	Token *Token
	Value Expression
}

func (stm ReturnStmt) CanBeStmt() {}

func (stm ReturnStmt) String() string { return "ReturnStmt.String()" }

func NewReturnStmtWithToken(tok *Token, expr Expression) *ReturnStmt {
	return &ReturnStmt{Token: tok, Value: expr}
}
