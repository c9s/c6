package ast

type Op struct {
	Type  TokenType
	Token *Token
}

func NewOpWithToken(token *Token) *Op {
	return &Op{token.Type, token}
}

func NewOp(opType TokenType, token *Token) *Op {
	return &Op{opType, token}
}

func (op Op) String() string {
	if op.Token != nil {
		return op.Token.Str
	}
	return string(op.Type)
}
