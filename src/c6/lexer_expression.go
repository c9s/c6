package c6

import "unicode"

func lexExpression(l *Lexer) stateFn {

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
		l.emit(T_MINUS)
		return lexExpression

	} else if r == '*' {

		l.next()
		l.emit(T_MUL)
		return lexExpression

	} else if r == '+' {

		l.next()
		l.emit(T_PLUS)
		return lexExpression

	} else if r == '/' {

		l.next()
		l.emit(T_DIV)
		return lexExpression

	} else if r == '$' {

		lexVariableName(l)
		return lexExpression

	}
	return nil
}
