package ast

type Statement interface{}

/*
The nested statement allows declaration block and statements
*/
type NestedStatement struct{}

type AtRuleImport struct {
	Url       interface{} // if it's wrapped with url(...) or "string"
	MediaList []string
}

// for Url()
type Url string
type RelativeUrl string
