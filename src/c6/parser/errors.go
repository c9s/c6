package parser

import "fmt"
import "c6/ast"

/*
User's fault, probably.

Struct for common syntax error.

Examples:


panic(SyntaxError{
	Reason: ...,
	ActualToken: tok,
})
*/
type SyntaxError struct {
	Reason      string
	ActualToken *ast.Token
	Guide       string
	GuideUrl    string
	File        *ast.File
	// TODO: provide correction later
}

func (err SyntaxError) Error() (out string) {
	out = "Syntax error"
	if err.ActualToken != nil {
		out += fmt.Sprintf(" at line %d, offset %d. got %s\n", err.ActualToken.Line, err.ActualToken.LineOffset, err.ActualToken.Type.String())
	}
	if err.File != nil {
		out += fmt.Sprintf(" in file %s\n", err.File)
	}

	if err.Reason != "" {
		out += err.Reason + "\n"
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
