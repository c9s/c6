package lexer

import (
	"github.com/c9s/c6/ast"
)

func lexCommentLine(l *Lexer) stateFn {
	if !l.match("//") {
		return nil
	}
	l.ignore()

	var r = l.next()
	for r != EOF {
		if r == '\n' {
			break
		}
		r = l.next()
		if r == '\r' {
			r = l.next()
		}
	}
	l.backup()
	l.ignore()
	return lexStart
}

func lexCommentBlock(l *Lexer, emit bool) stateFn {
	if !l.match("/*") {
		return nil
	}
	l.ignore()
	var r = l.next()
	for r != EOF {
		if r == '*' && l.peek() == '/' {
			l.backup()
			if emit {
				l.emit(ast.T_COMMENT_BLOCK)
			} else {
				l.ignore()
			}
			l.match("*/")
			l.ignore()
			return lexStart
		}
		r = l.next()
	}
	l.errorf("Expecting comment end mark '*/'. Got '%c'", r)
	return lexStart
}

func lexComment(l *Lexer, emit bool) stateFn {
	var r = l.peek()
	var r2 = l.peekBy(2)
	if r == '/' && r2 == '*' {
		lexCommentBlock(l, emit)
	} else if r == '/' && r2 == '/' {
		lexCommentLine(l)
	}
	return nil
}
