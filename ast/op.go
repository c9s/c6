package ast

type Op struct {
	Type  TokenType
	Token *Token
}

func NewOpWithToken(token *Token) *Op {
	return &Op{token.Type, token}
}

func NewOp(opType TokenType) *Op {
	return &Op{opType, nil}
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
	case T_PAREN_OPEN:
		return "("
	case T_PAREN_CLOSE:
		return ")"
	case T_NOP:
		return ""
	}
	panic("Unsupported token type")
	///XXX return ""
}

func (op Op) String() string {
	if op.Token != nil {
		return op.Token.Str
	}
	return OpTokenName(op.Type)
}
