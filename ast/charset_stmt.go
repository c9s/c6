package ast

type CharsetStmt struct {
	Encoding string
	Token    *Token
}

func (self CharsetStmt) CanBeStmt() {}

func (self CharsetStmt) String() string {
	return "@charset " + self.Encoding + ";"
}

func NewCharsetStmtWithToken(token *Token) *CharsetStmt {
	return &CharsetStmt{token.Str, token}
}
