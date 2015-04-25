package ast

type Op struct {
	Op    string
	Token *Token
}

func NewOp(token *Token) *Op {
	return &Op{token.Str, token}
}
