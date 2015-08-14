package ast

type MediaQueryStmt struct {
	MediaQueryList *MediaQueryList
	Block          *DeclBlock
	Scope          *Scope
}

type MediaQueryList []*MediaQuery

func (list MediaQueryList) Append(query *MediaQuery) {
	newlist := append(list, query)
	list = newlist
}

func NewMediaQueryList() *MediaQueryList {
	return &MediaQueryList{}
}

func (stm MediaQueryStmt) CanBeStmt() {}

func NewMediaQueryStmt() *MediaQueryStmt {
	return &MediaQueryStmt{}
}

func (stm MediaQueryStmt) String() (out string) {
	for _, mediaQuery := range *stm.MediaQueryList {
		out += ", " + mediaQuery.String()
	}
	return out[2:]
}

/*
One MediaQuery may contain media type or media expression.
*/
type MediaQuery struct {
	MediaType Expr
	MediaExpr Expr
}

func NewMediaQuery(mediaType Expr, expr Expr) *MediaQuery {
	return &MediaQuery{mediaType, expr}
}

func (stm MediaQuery) CSS3String() string {
	return stm.String()
}

func (stm MediaQuery) String() (out string) {
	/*
		{media type} and {media expression}
	*/
	if stm.MediaType != nil {
		out += stm.MediaType.String()
	}
	if stm.MediaExpr != nil {
		if stm.MediaType != nil {
			out += " and "
		}
		out += stm.MediaExpr.String()
	}
	return out
}
