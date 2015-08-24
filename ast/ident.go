package ast

type Ident struct {
	Ident string
	Token *Token
}

func (self Ident) CanBeNode() {}
func (self Ident) String() string {
	return self.Ident
}

func NewIdentWithToken(token *Token) *Ident {
	return &Ident{token.Str, token}
}
