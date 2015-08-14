package ast

type Function struct {
	Ident        *Token
	ArgumentList *ArgumentList
	Block        *Block
}

func (f Function) CanBeStmt()     {}
func (f Function) String() string { return "Function.String() is unimplemented." }

func NewFunctionWithToken(tok *Token) *Function {
	return &Function{Ident: tok}
}
