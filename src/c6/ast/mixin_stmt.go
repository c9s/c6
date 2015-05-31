package ast

// import ""

type MixinStmt struct {
	Token        *Token
	Ident        *Token
	Block        *DeclBlock
	ArgumentList *ArgumentList
	Scope        *Scope
}

func (stm MixinStmt) CanBeStmt()     {}
func (stm MixinStmt) String() string { return "{mixin}" }

func NewMixinStmtWithToken(tok *Token) *MixinStmt {
	return &MixinStmt{Token: tok}
}
