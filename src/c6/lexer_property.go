package c6

import "unicode"

// import "strings"
func lexPropertyName(l *Lexer) stateFn {
	var r rune = l.next()
	for r == '-' || unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_PROPERTY_NAME)
	lexColon(l)
	return lexPropertyValue
}

func lexColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r = l.next()
	if r == ':' {
		l.emit(T_COLON)
	} else {
		l.error("Expecting ':' token, Got '%s'", r)
	}
	return nil
}

func lexPropertyValue(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.peek()

	if r == '#' && l.peekMore(2) == '{' {
		return lexInterpolation
	} else if r == '#' {
		return lexHexColor
	} else if unicode.IsDigit(r) {
		return lexNumber
	} else if unicode.IsLetter(r) {

		l.remember()

		r = l.next()

		return lexConstantString

	} else if r == '-' {

		l.next()
		l.emit(T_MINUS)
		return lexPropertyValue

	} else if r == '+' {

		l.next()
		l.emit(T_PLUS)
		return lexPropertyValue

	} else if r == '/' {
		l.next()
		l.emit(T_DIV)
		return lexPropertyValue
	} else if r == '*' {
		l.next()
		l.emit(T_MUL)
		return lexPropertyValue
	} else if r == '$' {
		return lexVariableName
	} else if r == ' ' {
		l.next()
		l.ignore()
		return lexPropertyValue
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStatement
	} else if r == '}' {
		l.next()
		l.emit(T_BRACE_END)
		return lexStart
	} else if r == EOF {
		l.error("Unexpected end of file", r)
	} else {
		l.error("can't lex rune for property value: '%s'", r)
	}
	return nil
}
