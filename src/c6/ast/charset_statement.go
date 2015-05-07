package ast

type CharsetStatement struct {
	Charset string
	Token   *Token
}

func (self CharsetStatement) CanBeStatement() {}

func (self CharsetStatement) String() string {
	return "@charset " + self.Charset + ";"
}

func NewCharsetStatementWithToken(token *Token) *CharsetStatement {
	return &CharsetStatement{token.Str, token}
}
