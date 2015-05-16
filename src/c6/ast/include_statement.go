package ast

type IncludeStatement struct {
	Token        *Token // @include
	MixinIdent   *Token // mixin identitfier
	ArgumentList *ArgumentList
	ContentBlock *DeclarationBlock // if any
}

func (stm IncludeStatement) CanBeStatement() {}
func (stm IncludeStatement) String() string  { return "IncludeStatement.String()" }

func NewIncludeStatementWithToken(token *Token) *IncludeStatement {
	return &IncludeStatement{
		Token: token,
	}
}
