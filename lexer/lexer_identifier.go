package lexer

import "unicode"
import "github.com/c9s/c6/ast"

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
