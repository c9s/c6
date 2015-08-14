package ast

type ReturnStmt struct {
	Token *Token
	Value Expr
}

func (stm ReturnStmt) CanBeStmt() {}

func (stm ReturnStmt) String() string { return "ReturnStmt.String()" }

func NewReturnStmtWithToken(tok *Token, expr Expr) *ReturnStmt {
	return &ReturnStmt{Token: tok, Value: expr}
}
