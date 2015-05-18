package ast

type Url interface{}

/*
For relative Url
*/
type RelativeUrl string

/*
For url(http:....) or "http://...."
*/
type AbsoluteUrl string

/*
For import like this:

@import "component/list"; // => component/_list.scss
*/
type ScssImportUrl string

/**
The @import rule syntax is described here:

@see http://www.w3.org/TR/2015/CR-css-cascade-3-20150416/#at-import
*/
type ImportStatement struct {
	Url            Url // if it's wrapped with url(...) or "string"
	MediaQueryList []*MediaQuery
}

func NewImportStatement() *ImportStatement {
	return &ImportStatement{Url: nil}
}

func (self ImportStatement) CanBeStatement() {}

func (self ImportStatement) String() string { return "ImportStatement.String()" }
