package c6

import "unicode"
import "c6/ast"

func lexFunctionParams(l *Lexer) stateFn {
	var r = l.next()
	if r != '(' {
		l.error("Expecting '('. Got '%s'.", r)
	}
	l.emit(ast.T_PAREN_START)
	l.ignoreSpaces()

	r = l.peek()
	for r != EOF {
		l.ignoreSpaces()
		r = l.peek()
		if r == ')' {
			l.next()
			l.emit(ast.T_PAREN_END)
			break
		}
		if lexExpression(l) == nil {
			break
		}

		l.ignoreSpaces()
		r = l.peek()
		if r == ',' {
			l.next()
			l.emit(ast.T_COMMA)
		}
	}
	return nil
}

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

	if l.peek() == '(' {
		l.emit(ast.T_FUNCTION_NAME)
		lexFunctionParams(l)
	} else {
		l.emit(ast.T_IDENT)
	}
	return lexExpression
}

/*
Lexing expression with interpolation support.
*/
func lexExpression(l *Lexer) stateFn {
	var leadingSpaces = l.ignoreSpaces()

	var r = l.peek()
	var r2 = l.peekBy(2)
	var lastToken = l.lastToken()

	// avoid double literal concat
	if lastToken != nil && lastToken.Type != ast.T_LITERAL_CONCAT {
		if leadingSpaces == 0 && lastToken != nil && lastToken.Type == ast.T_INTERPOLATION_END {
			l.emit(ast.T_LITERAL_CONCAT)
		} else if leadingSpaces == 0 && l.Offset > 0 && r == '#' && r2 == '{' {
			l.emit(ast.T_LITERAL_CONCAT)
		}
	}

	if l.matchKeywordMap(exprTokenMap) {

	} else if r == 'U' && r2 == '+' {

		lexUnicodeRange(l)

	} else if unicode.IsLetter(r) {

		lexIdentifier(l)

	} else if r == '.' && unicode.IsDigit(r2) {

		// lexNumber may return lexNumber unit
		if fn := lexNumber(l); fn != nil {
			fn(l)
		}

	} else if unicode.IsDigit(r) {

		// lexNumber may return lexNumber unit
		if fn := lexNumber(l); fn != nil {
			fn(l)
		}

	} else if r == '-' {
		var r2 = l.peekBy(2)

		if unicode.IsLetter(r2) {
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

		if r2 == '*' {
			// don't emit the comment inside the expression
			lexComment(l, false)
		} else {
			l.next()
			l.emit(ast.T_DIV)
		}

	} else if r == ':' { // a port of map

		l.next()
		l.emit(ast.T_COLON)

	} else if r == ',' { // a part of map or list

		l.next()
		l.emit(ast.T_COMMA)

	} else if r == '(' {

		l.next()
		l.emit(ast.T_PAREN_START)

	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_END)

	} else if r == '=' {

		l.next()
		l.emit(ast.T_EQUAL)

	} else if r == '#' {

		// ignore interpolation here, we need to handle interpolation in the
		// caller because we need to know the context...  interpolation is the
		// tricky part we need to handle, we need to think about a better
		// solution here..
		if l.peekBy(2) == '{' {

			lexInterpolation2(l)

		} else {
			lexHexColor(l)
		}

	} else if r == '"' || r == '\'' {

		lexString(l)

	} else if r == '$' {

		lexVariableName(l)

	} else if r == EOF {

		return nil

	} else {

		// for ';' and '}'
		return nil

	}

	// for interpolation after any token above
	if l.peek() == '#' && l.peekBy(2) == '{' {
		l.emit(ast.T_LITERAL_CONCAT)
		lexInterpolation2(l)
	}

	// the default return stats
	return lexExpression
}
