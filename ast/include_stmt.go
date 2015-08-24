package ast

type IncludeStmt struct {
	Token        *Token // @include
	MixinIdent   *Token // mixin identitfier
	ArgumentList *ArgumentList
	ContentBlock *DeclBlock // if any
}

func (stm IncludeStmt) CanBeStmt()     {}
func (stm IncludeStmt) String() string { return "IncludeStmt.String()" }

func NewIncludeStmtWithToken(token *Token) *IncludeStmt {
	return &IncludeStmt{
		Token: token,
	}
}
