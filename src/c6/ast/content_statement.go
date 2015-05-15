package ast

type ContentStatement struct {
	Token *Token
}

func (stm ContentStatement) CanBeStatement() {}

func (stm ContentStatement) String() string {
	return stm.Token.String()
}

func NewContentStatementWithToken(tok *Token) *ContentStatement {
	return &ContentStatement{tok}
}
