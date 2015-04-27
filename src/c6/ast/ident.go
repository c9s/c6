package ast

type Ident struct {
	Ident string
	Token Token
}

func (self Ident) CanBeExpression() {}
func (self Ident) CanBeNode()       {}

func NewIdent(ident string, token Token) *Ident {
	return &Ident{ident, token}
}
