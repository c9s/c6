package c6

type SelectorOp int

const (
	SelectorOpAnd        SelectorOp = iota
	SelectorOpDescendant            // E F
	SelectorOpChild                 // E > F
)

type Selector struct{}

type SelectorExpression struct {
	Operator        SelectorOp
	LeftExpression  *Selector
	RightExpression *Selector
}

type Rule struct {
	Selectors SelectorExpression
}
