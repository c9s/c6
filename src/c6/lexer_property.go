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

func lexMicrosoftProgIdFunction(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
		r = l.next()
	}
	l.backup()

	// here starts the sproperty
	r = l.next()
	if r != '(' {
		l.error("Expecting '(' after the MS function name. Got %s", r)
	}
	l.emit(ast.T_PAREN_START)

	l.ignoreSpaces()
	r = l.next()
	for r != ')' {
		// lex function parameter name
		for unicode.IsLetter(r) || unicode.IsDigit(r) {
			r = l.next()
		}
		l.emit(ast.T_PARAM_NAME)
		l.accept("=")
		l.emit(ast.T_EQUAL)

		l.ignoreSpaces()
		r = l.next()
	}
	l.emit(ast.T_PAREN_END)
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

	// for IE filter syntax like:
	//    progid:DXImageTransform.Microsoft.MotionBlur(strength=13, direction=310)
	if l.match("progid:") {
		l.emit(ast.T_MS_PROGID)
		lexMicrosoftProgIdFunction(l)
	}

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

	// try ';' or '}'
	r = l.next()
	if r == ';' {
		l.emit(ast.T_SEMICOLON)
	} else if r == '}' {
		// emit another semicolon here for parser simplicity?
		l.emit(ast.T_BRACE_END)
	} else {
		l.backup()
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
	l.ignoreSpaces()
	return nil
}
