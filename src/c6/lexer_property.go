package c6

import "unicode"
import "c6/ast"

/*
Possible property value syntax:

   width: 10px;        // numeric
   width: 10px + 10px; // expression
   border: 1px #{solid} #000;   // interpolation
   color: rgba( 0, 0, 255, 1.0);  // rgba function
   width: auto;    // string constant

*/
func lexProperty(l *Lexer) stateFn {

	// accept all leading slash
	for l.accept("-") {
	}

	var r = l.next()
	for {
		if r == '#' && l.peek() == '{' {
			l.backup() // back to the position before '#'
			lexInterpolation2(l)
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			break
		}
		r = l.next()
	}
	l.backup()

	if l.precedeStartOffset() {
		l.emit(ast.T_PROPERTY_NAME)
	}

	lexColon(l)

	r = l.peek()
	for r != ';' && r != '}' {
		lexExpression(l)
		r = l.peek()
	}
	r = l.next()
	if r == ';' {
		l.emit(ast.T_SEMICOLON)
	} else if r == '}' {
		l.emit(ast.T_BRACE_END)
	}
	return lexStatement
}

func lexColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r = l.next()
	if r != ':' {
		l.error("Expecting ':' token, Got '%s'", r)
	}
	l.emit(ast.T_COLON)
	return nil
}
