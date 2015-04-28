package ast

type Op struct {
	Op    string
	Token *Token
}

func NewOp(token *Token) *Op {
	return &Op{token.Str, token}
}

func (self Op) String() string {
	return self.Op
}
