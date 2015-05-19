package ast

type Statement interface {
	CanBeStatement()
	String() string
}

type StatementList []Statement

func (list StatementList) Append(stm Statement) {
	newlist := append(list, stm)
	list = newlist
}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

func (stm NestedStatement) CanBeStatement() {}
