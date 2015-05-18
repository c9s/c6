package c6

import "fmt"
import "c6/ast"

/*
User's fault, probably.

Struct for common syntax error.

Examples:


panic(SyntaxError{
	Expecting: ...,
	ActualToken: tok,
})
*/
type SyntaxError struct {
	Expecting   string
	ActualToken *ast.Token
	Guide       string
	GuideUrl    string
	// TODO: provide correction later
}

func (err SyntaxError) Error() (out string) {
	out = "Syntax error "
	if err.ActualToken != nil {
		out += fmt.Sprintf(" at line %d, offset %d. given %s\n", err.ActualToken.Line, err.ActualToken.Pos, err.ActualToken.Type.String())
	}
	if err.Expecting != "" {
		out += "The parser expects " + err.Expecting + "\n"
	}
	if err.Guide != "" {
		out += "We suggest you to " + err.Guide + "\n"
	}
	if err.GuideUrl != "" {
		out += "For more information, please visit " + err.GuideUrl + "\n"
	}
	return out
}

/*
Parser's fault
*/
type ParserError struct {
	ExpectingToken string
	ActualToken    string
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Expecting '%s', but the actual token we got was '%s'.", e.ExpectingToken, e.ActualToken)
}
