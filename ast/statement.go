package ast

type Stmt interface {
	CanBeStmt()
	String() string
}

type StmtList []Stmt

func (list StmtList) Append(stm Stmt) {
	newlist := append(list, stm)
	list = newlist
}

/*
The nested statement allows declaration block and statements
*/
type NestedStmt struct{}

func (stm NestedStmt) CanBeStmt() {}
