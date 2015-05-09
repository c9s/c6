package ast

type MediaQueryStatement struct {
	MediaQueryList []*MediaQuery
}

func (stm MediaQueryStatement) CanBeStatement() {}

func NewMediaQueryStatement() *MediaQueryStatement {
	return &MediaQueryStatement{}
}

func (stm MediaQueryStatement) String() (out string) {
	for _, mediaQuery := range stm.MediaQueryList {
		out += ", " + mediaQuery.String()
	}
	return out[2:]
}

/*
One MediaQuery may contain media type or media expression.
*/
type MediaQuery struct {
	MediaType       Expression
	MediaExpression Expression
}

func NewMediaQuery(mediaType Expression, expr Expression) *MediaQuery {
	return &MediaQuery{mediaType, expr}
}

func (stm MediaQuery) String() (out string) {
	if stm.MediaType != nil {
		out += stm.MediaType.String()
	}
	if stm.MediaExpression != nil {
		if stm.MediaType != nil {
			out += " and "
		}
		out += stm.MediaExpression.String()
	}
	return out
}
