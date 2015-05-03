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

func lexExpression(l *Lexer) stateFn {
	l.ignoreSpaces()

	var r = l.peek()

	if r == 't' && l.match("true") {

		l.emit(ast.T_TRUE)

	} else if r == 'f' && l.match("false") {

		l.emit(ast.T_FALSE)
		return lexExpression

	} else if r == 'n' && l.match("null") {

		l.emit(ast.T_NULL)

	} else if r == 'a' && l.match("and") {

		l.emit(ast.T_AND)

	} else if r == 'o' && l.match("or") {

		l.emit(ast.T_OR)

	} else if r == 'x' && l.match("xor") {

		l.emit(ast.T_XOR)

	} else if unicode.IsDigit(r) {

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

		l.next()
		l.emit(ast.T_DIV)

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

		// ignore interpolation
		if l.peekBy(2) == '{' {
			return nil
		}

		lexHexColor(l)

	} else if r == '"' || r == '\'' {

		lexString(l)

	} else if r == '$' {

		lexVariableName(l)

	} else if unicode.IsLetter(r) {

		lexIdentifier(l)

	} else if r == EOF {

		return nil

	} else {

		// for ';' and '}'
		return nil

	}

	// the default return stats
	return lexExpression
}
