package lexer

import "unicode"
import "github.com/c9s/c6/ast"

func lexFunctionParams(l *Lexer) stateFn {
	var r = l.next()
	if r != '(' {
		l.errorf("Expecting '('. Got '%c'.", r)
	}
	l.emit(ast.T_PAREN_OPEN)
	l.ignoreSpaces()

	r = l.peek()
	for r != EOF {
		l.ignoreSpaces()
		r = l.peek()
		if r == ')' {
			l.next()
			l.emit(ast.T_PAREN_CLOSE)
			break
		}
		if lexExpr(l) == nil {
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
		panic("An identifier needs to start with a letter or dash")
	}
	r = l.next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		if r == '-' {
			var r2 = l.peek()
			if !unicode.IsLetter(r2) && r2 != '-' {
				l.backup()
				return lexExpr
			}
		}

		r = l.next()
	}
	l.backup()

	if l.peek() == '(' {
		var curTok = l.emit(ast.T_FUNCTION_NAME)

		if curTok.Str == "url" || curTok.Str == "local" {
			lexUrlParam(l)
		} else {
			lexFunctionParams(l)
		}
	} else {
		l.emit(ast.T_IDENT)
	}
	return lexExpr
}

/*
Lexing expression with interpolation support.
*/
func lexExpr(l *Lexer) stateFn {
	var leadingSpaces = l.ignoreSpaces()

	var r, r2 = l.peek2()
	var lastToken = l.lastToken()

	// avoid double literal concat
	if lastToken != nil && lastToken.Type != ast.T_LITERAL_CONCAT {
		if leadingSpaces == 0 && lastToken != nil && lastToken.Type == ast.T_INTERPOLATION_END {
			l.emit(ast.T_LITERAL_CONCAT)
		} else if leadingSpaces == 0 && l.Offset > 0 && r == '#' && r2 == '{' {
			l.emit(ast.T_LITERAL_CONCAT)
		}
	}

	if l.matchKeywordList(ast.ExprTokenList) != nil {

	} else if l.matchKeywordList(ast.FlagTokenList) != nil {

	} else if l.matchKeywordList(ast.ForRangeKeywordTokenList) != nil {

	} else if r == 'U' && r2 == '+' {

		lexUnicodeRange(l)

	} else if unicode.IsLetter(r) {

		lexIdentifier(l)

	} else if r == '.' && r2 == '.' {

		l.next()
		l.next()
		l.emit(ast.T_RANGE)

	} else if r == '.' && unicode.IsDigit(r2) {

		// lexNumber may return lexNumber unit
		if fn := lexNumber(l); fn != nil {
			fn(l)
		}

	} else if unicode.IsDigit(r) {

		if fn := lexNumber(l); fn != nil {
			fn(l)
		}

	} else if r == '+' {

		l.next()
		l.emit(ast.T_PLUS)

	} else if r == '-' {

		if r2 == 'n' && !unicode.IsLetter(l.peekBy(3)) {

			l.next()
			l.emit(ast.T_MINUS)

			l.accept("n")
			l.emit(ast.T_N)

		} else if unicode.IsLetter(r2) {
			// XXX: Works for '-moz' or '-webkit-..' but we should move this to property lexing...
			//    like:
			//         background-image: -moz-linear-gradient(top, #81a8cb, #4477a1);

			l.next()
			lexIdentifier(l)

		} else {

			l.next()
			l.emit(ast.T_MINUS)

		}

	} else if r == '*' {

		l.next()
		l.emit(ast.T_MUL)

	} else if r == '&' {

		// '&' is allowed in expression, to make sure there is a parent selector
		l.next()
		l.emit(ast.T_PARENT_SELECTOR)

	} else if r == '%' {

		// TODO: placeholders start with '%'
		l.next()
		l.emit(ast.T_MOD)

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
		l.emit(ast.T_PAREN_OPEN)

	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_CLOSE)

	} else if r == '<' {

		l.next()
		if r2 == '=' {
			l.next()
			l.emit(ast.T_LE)
		} else {
			l.emit(ast.T_LT)
		}

	} else if r == '>' {
		l.next()
		if r2 == '=' {
			l.next()
			l.emit(ast.T_GE)
		} else {
			l.emit(ast.T_GT)
		}

	} else if r == '!' && r2 == '=' {

		l.next()
		l.next()
		l.emit(ast.T_UNEQUAL)

	} else if r == '=' {

		l.next()

		if r2 == '=' {

			l.next()
			l.emit(ast.T_EQUAL)

		} else {

			l.emit(ast.T_ASSIGN)

		}

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

	} else if r == EOF || r == '}' || r == '{' || r == ';' { // let expression lexer stop before the start or end of block.

		return nil

	} else {

		// anything else expression lexer don't know
		return nil

	}

	// for interpolation after any token above
	r, r2 = l.peek2()
	if r == '#' && r2 == '{' {
		l.emit(ast.T_LITERAL_CONCAT)
		lexInterpolation2(l)
	}

	// the default return stats
	return lexExpr
}
