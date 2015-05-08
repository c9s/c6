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

func OpTokenName(tokType TokenType) string {
	switch tokType {
	case T_DIV:
		return "/"
	case T_MUL:
		return "*"
	case T_MINUS:
		return "-"
	case T_PLUS:
		return "+"
	case T_PAREN_START:
		return "("
	case T_PAREN_END:
		return ")"
	case T_NOP:
		return ""
	}
	panic("Unsupported token type")
	return ""
}

func (op Op) String() string {
	if op.Token != nil {
		return op.Token.Str
	}
	return OpTokenName(op.Type)
}
