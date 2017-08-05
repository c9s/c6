package lexer

import (
	"github.com/c9s/c6/ast"
	"unicode"
)

func lexStmt(l *Lexer) stateFn {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r, r2 = l.peek2()

	if r == EOF {
		return nil
	}

	if r == '@' {

		return lexAtRule

	} else if r == '(' {
		l.next()
		l.emit(ast.T_PAREN_OPEN)
		return lexStart
	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_CLOSE)
		return lexStart

	} else if r == '{' {

		l.next()
		l.emit(ast.T_BRACE_OPEN)
		return lexStmt

	} else if r == '}' {

		l.next()
		l.emit(ast.T_BRACE_CLOSE)
		return lexStmt

	} else if l.match("<!--") {

		l.emit(ast.T_CDO)
		return lexStmt

	} else if l.match("-->") {

		l.emit(ast.T_CDC)
		return lexStmt

	} else if r == '/' && (r2 == '*' || r2 == '/') {

		lexComment(l, true)

		return lexStmt

	} else if r == '$' { // it's a variable assignment statement

		return lexAssignStmt

	} else if r == ';' {

		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStart

	} else if r == ',' {

		l.next()
		l.emit(ast.T_COMMA)
		return lexStart

	} else if r == '-' {

		// lex the slash prefix property name
		return lexProperty

	} else if r == '#' && r2 != '{' {

		return lexSelectors

	} else if unicode.IsLetter(r) || (r == '#') { // it might be -vendor- property or a property name or a selector

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

	} else if looksLikeSelector(r) {

		return lexSelectors

	} else if r == EOF {

		return nil

	} else {

		l.errorf("Unexpected token: '%c'", r)

	}
	return nil
}

func looksLikeSelector(r rune) bool {
	return r == '[' || r == '*' || r == '>' || r == '&' || r == '#' || r == '.' || r == '+' || r == ':'
}
