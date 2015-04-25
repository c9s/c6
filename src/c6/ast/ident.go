package ast

type Ident struct {
	Ident string
	Token Token
}

func (self Ident) CanBeExpression() {}

func NewIdent(ident string, token Token) *Ident {
	return &Ident{ident, token}
}
