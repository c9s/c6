package lexer

import "unicode"
import "github.com/c9s/c6/ast"

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
	l.emit(ast.T_FUNCTION_NAME)

	// here starts the sproperty
	r = l.next()
	if r != '(' {
		l.errorf("Expecting '(' after the MS function name. Got %c", r)
	}
	l.emit(ast.T_PAREN_OPEN)

	l.ignoreSpaces()

	// here comes the sProperty
	//     progid:DXImageTransform.Microsoft.filtername(sProperties)
	// @see https://msdn.microsoft.com/en-us/library/ms532847(v=vs.85).aspx
	for r != ')' {
		// lex function parameter name
		r = l.next()
		for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			r = l.next()
		}
		l.backup()
		l.emit(ast.T_MS_PARAM_NAME)
		l.accept("=")
		l.emit(ast.T_ATTR_EQUAL)

		lexExpr(l)

		l.ignoreSpaces()
		r = l.peek()
		if r == ',' {
			l.next()
			l.emit(ast.T_COMMA)
			l.ignoreSpaces()
		} else if r == ')' {
			l.next()
			l.emit(ast.T_PAREN_CLOSE)
			break
		}
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
				l.emit(ast.T_LITERAL_CONCAT)
			}
		}

		// we have something
		if lexPropertyNameToken(l) != nil {
			r = l.peek()
			if !unicode.IsSpace(r) && r != ':' {
				l.emit(ast.T_LITERAL_CONCAT)
			}
		}
		r = l.peek()
	}

	lexColon(l)

	l.remember()
	l.ignoreSpaces()

	// for IE filter syntax like:
	//    progid:DXImageTransform.Microsoft.MotionBlur(strength=13, direction=310)
	if l.match("progid:") {
		l.emit(ast.T_MS_PROGID)
		lexMicrosoftProgIdFunction(l)
	} else {
		l.rollback()
	}

	// the '{' is used for the start token of nested properties
	r = l.peek()
	for r != EOF {

		// for nested property value
		if r == '{' {

			l.next()
			l.emit(ast.T_BRACE_OPEN)
			return lexStmt

		} else if lexExpr(l) == nil {
			break
		}

		r = l.peek()
	}

	l.ignoreSpaces()
	lexComment(l, false)
	l.ignoreSpaces()

	// the semicolon in the last declaration is optional.
	l.ignoreSpaces()
	if l.accept(";") {
		l.emit(ast.T_SEMICOLON)
	}

	l.ignoreSpaces()
	if l.accept("}") {
		l.emit(ast.T_BRACE_CLOSE)
	}
	return lexStmt
}

func lexColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r = l.next()
	if r != ':' {
		l.errorf("Expecting ':' token, Got '%c'", r)
	}
	l.emit(ast.T_COLON)

	// We don't ignore space after the colon because we need spaces to detect literal concat.
	return nil
}
