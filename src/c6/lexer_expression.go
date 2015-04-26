package c6

import "unicode"
import "c6/ast"

func lexIdentifier(l *Lexer) stateFn {
	var r = l.next()
	if !unicode.IsLetter(r) && r != '-' {
		panic("An identifier needs to start with a letter")
	}
	r = l.next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		if r == '-' {
			if !unicode.IsLetter(l.peek()) {
				l.backup()
				return lexExpression
			}
		}

		r = l.next()
	}
	l.backup()
	l.emit(ast.T_IDENT)
	return lexExpression
}

func lexExpression(l *Lexer) stateFn {
	l.ignoreSpaces()

	var r = l.peek()

	if r == 't' && l.match("true") {

		l.emit(ast.T_TRUE)

	} else if r == 'f' && l.match("false") {

		l.emit(ast.T_FALSE)
		return lexExpression

	} else if r == 'n' && l.match("null") {

		l.emit(ast.T_NULL)

	} else if unicode.IsDigit(r) {

		if fn := lexNumber(l); fn != nil {
			fn(l)
		}

	} else if r == '-' {

		if unicode.IsLetter(l.peekBy(2)) {
			lexIdentifier(l)
		} else {
			l.next()
			l.emit(ast.T_MINUS)
		}

	} else if r == '*' {

		l.next()
		l.emit(ast.T_MUL)

	} else if r == '+' {

		l.next()
		l.emit(ast.T_PLUS)

	} else if r == '/' {

		l.next()
		l.emit(ast.T_DIV)

	} else if r == '(' {

		l.next()
		l.emit(ast.T_PAREN_START)

	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_END)

	} else if r == ' ' {

		l.next()
		l.ignore()

	} else if r == ',' {

		l.next()
		l.emit(ast.T_COMMA)

	} else if r == '#' {

		if l.peekBy(2) == '{' {
			// It handles the concat after interpolation
			lexInterpolation2(l)
			return lexExpression

		} else {
			lexHexColor(l)
		}

	} else if r == '"' || r == '\'' {

		lexString(l)

	} else if r == '$' {

		lexVariableName(l)

	} else if unicode.IsLetter(r) {

		lexIdentifier(l)

	} else if r == EOF {

		return nil

	} else {

		return nil

	}

	if l.peek() == '#' && l.peekBy(2) == '{' {
		// inject a ast.T_CONCAT since we have an interpolation following the expression
		l.emit(ast.T_CONCAT)
	}

	// the default return stats
	return lexExpression
}
