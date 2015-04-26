package c6

import "c6/ast"
import "unicode"
import _ "fmt"

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
			for r == ' ' {
				r = l.next()
			}
			l.backup()

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
		l.error("Expecting interpolation token '#', Got %s", r)
	}
	r = l.next()
	if r != '{' {
		l.error("Expecting interpolation token '{', Got %s", r)
	}
	l.emit(ast.T_INTERPOLATION_START)

	// skip the space after #{
	l.ignoreSpaces()

	r = l.peek()
	for r != '}' {
		lexExpression(l)

		// ignore space
		l.ignoreSpaces()
		r = l.peek()

	}
	l.next() // consume '}'
	l.emit(ast.T_INTERPOLATION_END)

	r = l.peek()
	if !unicode.IsSpace(r) && r != '}' && r != ';' && r != ':' {
		l.emit(ast.T_CONCAT)
	}
	return nil
}
