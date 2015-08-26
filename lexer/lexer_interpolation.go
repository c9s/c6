package lexer

import "github.com/c9s/c6/ast"

import "unicode"

/*
There are 3 scope that users may use interpolation syntax:

 {selector interpolation}  {
	 {property name inpterpolation}: {property value interpolation}
 }

*/
func lexInterpolation(l *Lexer, emit bool) stateFn {
	l.remember()
	var r rune = l.next()
	if r == '#' {
		r = l.next()
		if r == '{' {
			if emit {
				l.emit(ast.T_INTERPOLATION_START)
			}

			r = l.next()
			for unicode.IsSpace(r) {
				r = l.next()
			}
			l.backup()

			// find the end of interpolation end brace
			r = l.next()
			for r != '}' {
				r = l.next()
			}
			l.backup()

			if emit {
				l.emit(ast.T_INTERPOLATION_INNER)
			}

			l.next() // for '}'
			if emit {
				l.emit(ast.T_INTERPOLATION_END)
			}
			return nil
		}
	}
	l.rollback()
	return nil
}

// Lex the expression inside interpolation
func lexInterpolation2(l *Lexer) stateFn {
	var r rune = l.next()
	if r != '#' {
		l.errorf("Expecting interpolation token '#', Got %c", r)
	}
	r = l.next()
	if r != '{' {
		l.errorf("Expecting interpolation token '{', Got %c", r)
	}
	l.emit(ast.T_INTERPOLATION_START)

	// skip the space after #{
	for lexExpr(l) != nil {
	}
	l.ignoreSpaces()
	l.expect("}")
	l.emit(ast.T_INTERPOLATION_END)
	return nil
}
