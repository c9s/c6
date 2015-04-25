package c6

import "unicode"
import "c6/ast"

func lexIdentifier(l *Lexer) stateFn {
	var r = l.next()
	if !unicode.IsLetter(r) {
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

	// ignore spaces
	var r = l.next()
	for r == ' ' {
		r = l.next()
	}
	l.backup()
	l.ignore()

	r = l.peek()

	if unicode.IsDigit(r) {

		if fn := lexNumber(l); fn != nil {
			fn(l)
		}
		return lexExpression

	} else if r == '-' {

		l.next()
		l.emit(ast.T_MINUS)
		return lexExpression

	} else if r == '*' {

		l.next()
		l.emit(ast.T_MUL)
		return lexExpression

	} else if r == '+' {

		l.next()
		l.emit(ast.T_PLUS)
		return lexExpression

	} else if r == '/' {

		l.next()
		l.emit(ast.T_DIV)
		return lexExpression

	} else if r == '(' {

		l.next()
		l.emit(ast.T_PAREN_START)
		return lexExpression

	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_END)
		return lexExpression

	} else if r == ' ' {

		l.next()
		l.ignore()
		return lexExpression

	} else if r == ',' {

		l.next()
		l.emit(ast.T_COMMA)
		return lexExpression

	} else if r == '#' {

		if l.peekBy(2) == '{' {
			lexInterpolation(l, true)
			return lexExpression
		}

		lexHexColor(l)
		return lexExpression

	} else if r == '$' {

		lexVariableName(l)
		return lexExpression

	} else if unicode.IsLetter(r) {

		lexIdentifier(l)
		return lexExpression

	} else if r == ';' {

		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStatement

	} else if r == '}' {

		l.next()
		l.emit(ast.T_BRACE_END)
		return lexStatement

	} else if r == EOF {
		// panic("Unexpected end of file")
		return nil
	}
	return nil
}
