package ast

// type Value struct { }
type Value interface {
	CanBeValue()
}
