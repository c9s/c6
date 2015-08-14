package ast

/**
The basic CSS compliant
*/
type CSSCompliant interface {
	CSSString() string
}

type CSS3Compliant interface {
	/*
		Output CSS3 compliant syntax
	*/
	CSS3String() string
}

type CSS4Compliant interface {
	/*
		Output CSS3 compliant syntax
	*/
	CSS4String() string
}
