package ast

type Statement interface {
	IsStatement()
}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

func (stm NestedStatement) IsStatement() {}

type ImportStatement struct {
	Url       interface{} // if it's wrapped with url(...) or "string"
	MediaList []string
}

func (self ImportStatement) IsStatement() {}

// for Url()
type Url string
type RelativeUrl string
