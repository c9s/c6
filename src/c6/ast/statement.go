package ast

type Statement interface {
	CanBeStatement()
}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

func (stm NestedStatement) CanBeStatement() {}
