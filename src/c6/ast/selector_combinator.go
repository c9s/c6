package ast

type Combinator interface {
	String() string
}

/*
Selectors present: E '+' F
*/
type AdjacentCombinator struct {
	Token *Token
}

func (self AdjacentCombinator) String() string    { return " + " }
func (self AdjacentCombinator) CSSString() string { return " + " }

func NewAdjacentCombinatorWithToken(token *Token) *AdjacentCombinator {
	return &AdjacentCombinator{token}
}

type DescendantCombinator struct {
	Token *Token
}

func NewDescendantCombinatorWithToken(token *Token) *DescendantCombinator {
	return &DescendantCombinator{token}
}

func NewDescendantCombinator() *DescendantCombinator {
	return &DescendantCombinator{}
}

func (self DescendantCombinator) String() string    { return " " }
func (self DescendantCombinator) CSSString() string { return " " }

type GeneralSiblingCombinator struct{ Token *Token }

func NewGeneralSiblingCombinator() *GeneralSiblingCombinator {
	return &GeneralSiblingCombinator{}
}

func NewGeneralSiblingCombinatorWithToken(token *Token) *GeneralSiblingCombinator {
	return &GeneralSiblingCombinator{token}
}

func (self GeneralSiblingCombinator) String() string    { return " ~ " }
func (self GeneralSiblingCombinator) CSSString() string { return " ~ " }

type ChildCombinator struct {
	Token *Token
}

func (self ChildCombinator) String() string    { return " > " }
func (self ChildCombinator) CSSString() string { return " > " }

func NewChildCombinatorWithToken(token *Token) *ChildCombinator {
	return &ChildCombinator{token}
}

func NewChildCombinator() *ChildCombinator {
	return &ChildCombinator{}
}
