package ast

type CharsetStatement struct {
	Encoding string
	Token    *Token
}

func (self CharsetStatement) CanBeStatement() {}

func (self CharsetStatement) String() string {
	return "@charset " + self.Encoding + ";"
}

func NewCharsetStatementWithToken(token *Token) *CharsetStatement {
	return &CharsetStatement{token.Str, token}
}
