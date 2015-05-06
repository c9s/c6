package ast

type ImportStatement struct {
	Url       interface{} // if it's wrapped with url(...) or "string"
	MediaList []string
}

func NewImportStatement() *ImportStatement {
	return &ImportStatement{nil, []string{}}
}

func (self ImportStatement) CanBeStatement() {}

func (self ImportStatement) String() string {
	return "@import ..."
}

// for Url()
type Url string

type RelativeUrl string
