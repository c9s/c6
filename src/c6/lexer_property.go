package c6

import "unicode"
import "c6/ast"

// import "strings"
func lexPropertyName(l *Lexer) stateFn {
	var r rune = l.next()
	for r == '-' || unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_PROPERTY_NAME)
	lexColon(l)
	return lexPropertyValue
}

func lexColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r = l.next()
	if r == ':' {
		l.emit(ast.T_COLON)
	} else {
		l.error("Expecting ':' token, Got '%s'", r)
	}
	return nil
}

/*
Possible property value syntax:

   width: 10px;        // numeric
   width: 10px + 10px; // expression
   border: 1px #{solid} #000;   // interpolation
   color: rgba( 0, 0, 255, 1.0);  // rgba function
   width: auto;    // string constant

*/
func lexPropertyValue(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.peek()

	if r == '#' {
		if l.peekBy(2) == '{' {
			lexInterpolation(l, true)
			return lexPropertyValue
		}
		return lexHexColor
	} else if unicode.IsDigit(r) {
		return lexNumber
	} else if unicode.IsLetter(r) {

		l.remember()

		r = l.next()

		return lexConstantString

	} else if r == '-' {

		l.next()
		l.emit(ast.T_MINUS)
		return lexPropertyValue

	} else if r == '+' {

		l.next()
		l.emit(ast.T_PLUS)
		return lexPropertyValue

	} else if r == '/' {
		l.next()
		l.emit(ast.T_DIV)
		return lexPropertyValue
	} else if r == '*' {
		l.next()
		l.emit(ast.T_MUL)
		return lexPropertyValue
	} else if r == ',' {
		l.next()
		l.emit(ast.T_COMMA)
		return lexPropertyValue
	} else if r == '$' {
		return lexVariableName
	} else if r == ' ' {
		l.next()
		l.ignore()
		return lexPropertyValue
	} else if r == ';' {
		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStatement
	} else if r == '}' {
		l.next()
		l.emit(ast.T_BRACE_END)
		return lexStart
	} else if r == EOF {
		l.error("Unexpected end of file", r)
	} else {
		l.error("can't lex rune for property value: '%s'", r)
	}
	return nil
}
