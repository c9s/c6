package ast

type Statement interface {
	IsStatement()
}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

func (stm NestedStatement) IsStatement() {}
