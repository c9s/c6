package c6

import "unicode"
import "c6/ast"

func lexPropertyNameToken(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		r = l.next()
	}
	l.backup()
	if l.precedeStartOffset() {
		l.emit(ast.T_PROPERTY_NAME_TOKEN)
		return lexPropertyNameToken
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
func lexProperty(l *Lexer) stateFn {
	var r = l.peek()
	for r != ':' {
		if l.peek() == '#' && l.peekBy(2) == '{' {
			lexInterpolation2(l)

			r = l.peek()
			if !unicode.IsSpace(r) && r != ':' {
				l.emit(ast.T_CONCAT)
			}
		}

		// we have something
		if lexPropertyNameToken(l) != nil {
			r = l.peek()
			if !unicode.IsSpace(r) && r != ':' {
				l.emit(ast.T_CONCAT)
			}
		}
		r = l.peek()
	}

	lexColon(l)

	l.ignoreSpaces()

	r = l.peek()
	for r != ';' && r != '}' {

		if l.peek() == '#' && l.peekBy(2) == '{' {
			lexInterpolation2(l)

			// See if it's the end of property
			r = l.peek()
			if !unicode.IsSpace(r) && r != '}' && r != ';' && r != ':' {
				l.emit(ast.T_CONCAT)
			}
		}

		if lexExpression(l) != nil {
			if l.peek() == '#' && l.peekBy(2) == '{' {
				l.emit(ast.T_CONCAT)
			}
		}
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
