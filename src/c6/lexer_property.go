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
func lexPropertyName(l *Lexer) stateFn {
	var r rune = l.next()
	for r == '-' || unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_PROPERTY_NAME)
	lexColon(l)
	return lexExpression
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
