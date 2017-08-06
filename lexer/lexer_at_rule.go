package lexer

import (
	"fmt"
	"github.com/c9s/c6/ast"
	"unicode"
)

/*
Currently the @import rule only supports '@import url(...) media;

@see https://developer.mozilla.org/en-US/docs/Web/CSS/@import for more @import syntax support
*/
func lexAtRule(l *Lexer) stateFn {
	var tok = l.matchKeywordList(ast.KeywordList)
	if tok != nil {
		switch tok.Type {
		case ast.T_IMPORT:
			l.ignoreSpaces()
			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_PAGE:
			l.ignoreSpaces()

			// lex pseudo selector ... if any
			if l.peek() == ':' {
				lexPseudoSelector(l)
			}
			return lexStart

		case ast.T_MEDIA:
			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_CHARSET:
			l.ignoreSpaces()
			return lexStart

		case ast.T_IF:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_ELSE_IF:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_ELSE:

			return lexStart

		case ast.T_FOR:

			return lexForStmt

		case ast.T_WHILE:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_CONTENT:
			return lexStart

		case ast.T_EXTEND:
			return lexSelectors

		case ast.T_FUNCTION, ast.T_RETURN, ast.T_MIXIN, ast.T_INCLUDE:
			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStart

		case ast.T_FONT_FACE:
			return lexStart

		default:
			var r = l.next()
			for unicode.IsLetter(r) {
				r = l.next()
			}
			l.backup()
			panic(fmt.Errorf("Unsupported at-rule directive '%s' %s", l.current(), tok))
		}
	}
	return nil
}
