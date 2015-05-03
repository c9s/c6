package ast

//go:generate stringer -type=OpType op.go token.go
type OpType int

const (
	OpNone OpType = iota
	OpAdd
	OpSub
	OpDiv
	OpMul
	OpPow
	OpConcat
	OpLiteralConcat
)

func ConvertTokenTypeToOpType(tokenType TokenType) OpType {
	switch tokenType {
	case T_PLUS:
		return OpAdd
	case T_MINUS:
		return OpSub
	case T_MUL:
		return OpMul
	case T_DIV:
		return OpDiv
	case T_LITERAL_CONCAT:
		return OpConcat
		// return OpLiteralConcat
	}
	panic("unknown token type")
	return OpNone
}

type Op struct {
	Type  OpType
	Token *Token
}

func NewOp(opType OpType, token *Token) *Op {
	return &Op{opType, token}
}
