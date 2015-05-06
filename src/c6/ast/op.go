package ast

//go:generate stringer -type=OpType op.go token.go
type OpType int

const (
	OpNone          OpType = 0
	OpAdd                  = T_PLUS
	OpSub                  = T_MINUS
	OpDiv                  = T_DIV
	OpMul                  = T_MUL
	OpConcat               = T_CONCAT
	OpLiteralConcat        = T_LITERAL_CONCAT
)

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
