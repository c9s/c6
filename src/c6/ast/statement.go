package ast

type Statement interface {
	CanBeStatement()
	String() string
}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

func (stm NestedStatement) CanBeStatement() {}
