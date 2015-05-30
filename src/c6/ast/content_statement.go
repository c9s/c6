package ast

type ContentStmt struct {
	Token *Token
}

func (stm ContentStmt) CanBeStmt() {}

func (stm ContentStmt) String() string {
	return stm.Token.String()
}

func NewContentStmtWithToken(tok *Token) *ContentStmt {
	return &ContentStmt{tok}
}
