package lexer

import (
	"github.com/c9s/c6/ast"
	"unicode"
)

func lexStart(l *Lexer) stateFn {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r, r2 = l.peek2()

	// lex simple statements
	switch r {

	case EOF:
		return nil

	case '@':
		return lexAtRule
	case '(':
		l.next()
		l.emit(ast.T_PAREN_OPEN)
		return lexStart
	case ')':
		l.next()
		l.emit(ast.T_PAREN_CLOSE)
		return lexStart

	case '{':
		l.next()
		l.emit(ast.T_BRACE_OPEN)
		return lexStart

	case '}':
		l.next()
		l.emit(ast.T_BRACE_CLOSE)
		return lexStart

	case '$':
		return lexAssignStmt

	case ';':
		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStart

	case '-':
		// Vendor prefix properties start with '-'
		return lexProperty

	case ',':
		l.next()
		l.emit(ast.T_COMMA)
		return lexStart

	case '/':
		if r2 == '*' {

			lexCommentBlock(l, true)
			return lexStart

		} else if r2 == '/' {

			return lexCommentLine

		}

		panic("unexpected token. expecing '*' or '/'")

	case '#':
		// make sure it's not an interpolation "#{" token
		if r2 != '{' {
			// looks like a selector
			return lexSelectors
		}

	case '[', '*', '>', '&', '.', '+', ':':
		return lexSelectors

	}

	if l.match("<!--") {

		l.emit(ast.T_CDOPEN)

		return lexStart
	}

	if l.match("-->") {

		l.emit(ast.T_CDCLOSE)

		return lexStart
	}

	// If a line starts with a letter or a sharp,
	// it might be a property name or selector, e.g.,
	//
	//    ul { }
	//
	//    -webkit-border-radius: 3px;
	//
	if unicode.IsLetter(r) || (r == '#') { // it might be -vendor- property or a property name or a selector

		// detect selector syntax
		l.remember()

		isProperty := false

		r = l.next()
		for r != EOF {
			// skip interpolation
			if r == '#' {
				if l.peek() == '{' {
					// find the matching brace
					r = l.next()
					for r != '}' {
						r = l.next()
					}
				}

			} else if r == ':' { // pseudo selector -> letters following ':', if there is a space after the ':' then it's a property value.

				if unicode.IsSpace(l.peek()) {
					isProperty = true
					break
				}

			} else if r == ';' {
				break
			} else if r == '}' {
				isProperty = true
				break
			} else if r == '{' {
				break
			} else if r == EOF {
				panic("unexpected EOF")
			}
			r = l.next()
		}

		l.rollback()

		if isProperty {
			return lexProperty
		} else {
			return lexSelectors
		}

	} else {

		l.errorf("Unexpected token: '%c'", r)

	}
	return nil
}
