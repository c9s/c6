package ast

type Op struct {
	Op    string
	Token Token
}

func NewOp(op string, token Token) *Op {
	return &Op{op, token}
}
