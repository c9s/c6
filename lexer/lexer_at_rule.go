package lexer

import (
	"fmt"
	"unicode"

	"github.com/c9s/c6/ast"
)

/*
Currently the @import rule only supports '@import url(...) media;

@see https://developer.mozilla.org/en-US/docs/Web/CSS/@import for more @import syntax support
*/
func lexAtRule(l *Lexer) (stateFn, error) {
	var tok = l.matchKeywordList(ast.KeywordList, false)
	if tok != nil {
		switch tok.Type {
		case ast.T_IMPORT:
			l.ignoreSpaces()
			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_PAGE:
			l.ignoreSpaces()

			// lex pseudo selector ... if any
			if l.peek() == ':' {
				if _, err := lexPseudoSelector(l); err != nil {
					return nil, err
				}
			}
			return lexStart, nil

		case ast.T_MEDIA:
			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_CHARSET:
			l.ignoreSpaces()
			return lexStart, nil

		case ast.T_DEBUG:
			fallthrough
		case ast.T_WARN:
			fallthrough
		case ast.T_ERROR:
			l.ignoreSpaces()
			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_IF:

			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_ELSE_IF:

			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_ELSE:

			return lexStart, nil

		case ast.T_FOR:

			return lexForStmt, nil

		case ast.T_WHILE:
			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_CONTENT:
			return lexStart, nil

		case ast.T_EXTEND:
			return lexSelectors, nil

		case ast.T_FUNCTION, ast.T_RETURN, ast.T_MIXIN, ast.T_INCLUDE:
			for {
				fn, err := lexExpr(l)
				if err != nil {
					return nil, err
				}

				if fn == nil {
					break
				}
			}
			return lexStart, nil

		case ast.T_FONT_FACE:
			return lexStart, nil

		case ast.T_AT_ROOT:
			return lexStart, nil

		default:
			var r = l.next()
			for unicode.IsLetter(r) {
				r = l.next()
			}
			l.backup()
			return nil, fmt.Errorf("Unsupported at-rule directive '%s' %s", l.current(), tok)
		}
	}
	return nil, nil
}
