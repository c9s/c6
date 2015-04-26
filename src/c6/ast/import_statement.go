package ast

type ImportStatement struct {
	Url       interface{} // if it's wrapped with url(...) or "string"
	MediaList []string
}

func (self ImportStatement) IsStatement() {}

// for Url()
type Url string
type RelativeUrl string
