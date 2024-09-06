package ast

type Stmt interface {
	CanBeStmt()
	String() string
}

type StmtList struct {
	Stmts []Stmt
}

func (list *StmtList) Append(stm ...Stmt) {
	list.Stmts = append(list.Stmts, stm...)
}

func (list *StmtList) AppendList(stm *StmtList) {
	if stm == nil {
		return
	}

	list.Stmts = append(list.Stmts, stm.Stmts...)
}

/*
The nested statement allows declaration block and statements
*/
type NestedStmt struct{}

func (stm NestedStmt) CanBeStmt() {}
