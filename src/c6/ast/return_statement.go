package ast

type ReturnStatement struct {
	Token *Token
	Value Expression
}

func (stm ReturnStatement) CanBeStatement() {}

func (stm ReturnStatement) String() string { return "ReturnStatement.String()" }

func NewReturnStatementWithToken(tok *Token, expr Expression) *ReturnStatement {
	return &ReturnStatement{Token: tok, Value: expr}
}
