package lexer

import (
	"fmt"
	"unicode"

	"github.com/c9s/c6/ast"
)

type stateFn func(*Lexer) stateFn

const LETTERS = "zxcvbnmasdfghjklqwertyuiop"
const DIGITS = "1234567890"

func (l *Lexer) errorf(msg string, r rune) {
	var err = fmt.Errorf(msg, string(r))
	panic(err)
}

func lexCommentLine(l *Lexer, emit bool) stateFn {
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
	if emit {
		l.emit(ast.T_COMMENT_LINE)
	} else {
		l.ignore()
	}
	return lexStmt
}

/*
Lex unicode range, used in `content` property.

@see https://developer.mozilla.org/en-US/docs/Web/CSS/unicode-range

Formal syntax: <urange>#
        where: <urange> = single_codepoint | codepoint_range | wildcard_range

	unicode-range: U+26               // single_codepoint
	unicode-range: U+0025-00FF        // codepoint_range
	unicode-range: U+4??              // wildcard_range
	unicode-range: U+0025-00FF, U+4?? // multiple values can be separated by commas

*/
func lexUnicodeRange(l *Lexer) stateFn {
	l.match("U+")

	var r = l.next()
	for unicode.IsDigit(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f') || r == '-' {
		r = l.next()
	}
	l.backup()

	if l.length() < 4 {
		panic(fmt.Errorf("Unicode-range requires at least 4 characters, we got %d. see https://developer.mozilla.org/en-US/docs/Web/CSS/unicode-range for more information", l.length()))
	}
	l.emit(ast.T_UNICODE_RANGE)
	return nil
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
			return lexStmt
		}
		r = l.next()
	}
	l.errorf("Expecting comment end mark '*/'. Got '%c'", r)
	return lexStmt
}

func lexComment(l *Lexer, emit bool) stateFn {
	var r = l.peek()
	var r2 = l.peekBy(2)
	if r == '/' && r2 == '*' {
		lexCommentBlock(l, emit)
	} else if r == '/' && r2 == '/' {
		lexCommentLine(l, emit)
	}
	return nil
}

func lexString(l *Lexer) stateFn {
	var r = l.next()
	if r == '"' {
		var containsInterpolation = false
		l.ignore()
		// string start
		r = l.next()
		for {
			if r == '"' {
				l.backup()
				token := l.createToken(ast.T_QQ_STRING)
				token.ContainsInterpolation = containsInterpolation
				l.emitToken(token)
				l.next()
				l.ignore()
				return lexStart
			} else if r == '\\' {
				// skip the escape character
			} else if IsInterpolationStartToken(r, l.peek()) {
				l.backup()
				lexInterpolation(l, false)
				containsInterpolation = true
			} else if r == EOF {
				panic("Expecting end of string")
			}
			r = l.next()
		}
		//XXX l.backup()
		//XXX return lexStart

	} else if r == '\'' {
		l.ignore()
		l.next()
		for {
			r = l.next()
			if r == '\'' {
				l.backup()
				l.emit(ast.T_Q_STRING)
				l.next()
				l.ignore()
				return lexStart
			} else if r == '\\' {
				// skip the escape character
				l.next()
			} else if r == EOF {
				panic("Expecting end of string")
			}
		}
		//XXX return lexStart
	}
	l.backup()
	return nil
}

func lexUrlParam(l *Lexer) {
	l.match("(")
	l.emit(ast.T_PAREN_OPEN)
	l.ignoreSpaces()

	var q = l.peek()
	if q == '"' || q == '\'' {
		lexString(l)
	} else {
		lexUnquoteStringExclude(l, "()")
	}

	l.ignoreSpaces()
	l.expect(")")
	l.emit(ast.T_PAREN_CLOSE)
}

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
			return lexStmt

		case ast.T_PAGE:
			l.ignoreSpaces()

			// lex pseudo selector ... if any
			if l.peek() == ':' {
				lexPseudoSelector(l)
			}
			return lexStmt

		case ast.T_MEDIA:
			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStmt

		case ast.T_CHARSET:
			l.ignoreSpaces()
			return lexStmt

		case ast.T_IF:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStmt

		case ast.T_ELSE_IF:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStmt

		case ast.T_ELSE:

			return lexStmt

		case ast.T_FOR:

			return lexForStmt

		case ast.T_WHILE:

			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStmt

		case ast.T_CONTENT:
			return lexStmt

		case ast.T_EXTEND:
			return lexSelectors

		case ast.T_FUNCTION, ast.T_RETURN, ast.T_MIXIN, ast.T_INCLUDE:
			for fn := lexExpr(l); fn != nil; fn = lexExpr(l) {
			}
			return lexStmt

		case ast.T_FONT_FACE:
			return lexStmt

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

func lexSpaces(l *Lexer) stateFn {
	for {
		var t = l.next()
		if t != ' ' {
			l.backup()
			return nil
		}
	}
	///XXX return lexStart
}

/*
lex unquote string but stops at the exclude rune.
*/
func lexUnquoteStringExclude(l *Lexer, exclude string) stateFn {
	l.til(exclude)
	l.emit(ast.T_UNQUOTE_STRING)
	return nil
}

func lexUnquoteStringStopAt(l *Lexer, stop rune) stateFn {
	var r = l.next()
	for r != stop {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_UNQUOTE_STRING)
	return nil
}

func lexUnquoteString(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_UNQUOTE_STRING)
	return nil
}

func lexAssignStmt(l *Lexer) stateFn {
	lexVariableName(l)
	lexColon(l)
	var r = l.peek()
	for r != ';' && r != '}' && r != EOF && lexExpr(l) != nil {
		r = l.peek()
	}
	// l.backup()

	l.ignoreSpaces()
	lexComment(l, false)
	l.ignoreSpaces()

	if l.accept(";") {
		l.emit(ast.T_SEMICOLON)
	} else if l.accept("}") {
		l.emit(ast.T_BRACE_CLOSE)
	}
	return lexStmt
}

func lexForStmt(l *Lexer) stateFn {
	l.ignoreSpaces()
	lexVariableName(l)

	fn := lexExpr(l)
	if fn == nil {
		panic("Expecting range expression after 'from'.")
	}
	for fn != nil {
		fn = lexExpr(l)
	}
	return lexStmt
}

// $var-rgba(255,255,0)
func lexVariableName(l *Lexer) stateFn {
	var r = l.next()
	if r != '$' {
		l.errorf("Unexpected token %c for lexVariable", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.errorf("The first character of a variable name must be letter. Got '%c'", r)
	}

	r = l.next()
	for r != EOF {
		if r == '-' {
			var r2 = l.peek()
			if unicode.IsLetter(r2) { // $a-b is a valid variable name.
				l.next()
			} else if unicode.IsDigit(r2) { // $a-3 should be $a '-' 3
				l.backup()
				l.emit(ast.T_VARIABLE)
				return lexExpr
			} else {
				break
			}
		} else if r == ':' {
			break
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		} else if r == '}' {
			l.backup()
			l.emit(ast.T_VARIABLE)
			return lexStmt
			///XXX break
		} else if unicode.IsSpace(r) || r == ';' {
			break
		}
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_VARIABLE)

	if l.match("...") {
		l.emit(ast.T_VARIABLE_LENGTH_ARGUMENTS)
	}

	return lexStmt
}

func lexHexColor(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.next()
	if r != '#' {
		l.errorf("Expecting hex color, got '%c'", r)
	}

	r = l.next()
	for unicode.In(r, unicode.ASCII_Hex_Digit) {
		r = l.next()
	}
	l.backup()

	var length = l.length() - 1
	if length != 3 && length != 6 && length != 8 {
		panic(fmt.Errorf("Invalid hex color, expecting 3, 6 or 8 hex characters, got %d - %s", length, l.current()))
	}
	l.emit(ast.T_HEX_COLOR)
	return lexExpr
}

/**
CSS time unit

@see https://developer.mozilla.org/zh-TW/docs/Web/CSS/time
*/
func lexNumberUnit(l *Lexer) stateFn {
	tok := l.matchKeywordList(ast.UnitTokenList)

	if tok == nil {
		var r = l.next()

		// for an+b syntax
		if r == 'n' && !unicode.IsLetter(l.peek()) {

			l.emit(ast.T_N)

		} else {

			// for other unit tokens
			for unicode.IsLetter(r) {
				r = l.next()
			}
			l.backup()
			if l.length() > 0 {
				l.emit(ast.T_UNIT_OTHERS)
			}

		}
	}

	if l.peek() == ';' {
		return lexStmt
	}
	return lexExpr
}

/**
@see https://developer.mozilla.org/en-US/docs/Web/CSS/number
*/
func lexNumber(l *Lexer) stateFn {
	var r = l.next()

	var floatPoint = false

	// allow floating number started with '.'
	if r == '.' {
		r = l.next()
		if !unicode.IsDigit(r) {
			l.errorf("Expecting digits after '.'. Got %c", r)
		}
		floatPoint = true
	}

	for unicode.IsDigit(r) {
		r = l.next()
		if r == '.' {
			floatPoint = true
			r = l.next()
			if !unicode.IsDigit(r) {
				l.errorf("Expecting at least one digit after the floating point, got '%c'", r)
			}
		} else if r == 'e' {
			var r2, r3 = l.peek2()
			// not scientific notation
			if !unicode.IsDigit(r2) && (r2 != '-' && !unicode.IsDigit(r3)) {
				break
			}
		}
	}
	l.backup()

	if floatPoint {
		l.emit(ast.T_FLOAT)
	} else {
		l.emit(ast.T_INTEGER)
	}
	return lexNumberUnit
}

func lexStmt(l *Lexer) stateFn {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r, r2 = l.peek2()

	if r == EOF {
		return nil
	}

	if r == '@' {

		return lexAtRule

	} else if r == '(' {
		l.next()
		l.emit(ast.T_PAREN_OPEN)
		return lexStart
	} else if r == ')' {

		l.next()
		l.emit(ast.T_PAREN_CLOSE)
		return lexStart

	} else if r == '{' {

		l.next()
		l.emit(ast.T_BRACE_OPEN)
		return lexStmt

	} else if r == '}' {

		l.next()
		l.emit(ast.T_BRACE_CLOSE)
		return lexStmt

	} else if l.match("<!--") {

		l.emit(ast.T_CDO)
		return lexStmt

	} else if l.match("-->") {

		l.emit(ast.T_CDC)
		return lexStmt

	} else if r == '/' && (r2 == '*' || r2 == '/') {

		lexComment(l, true)

		return lexStmt

	} else if r == '$' { // it's a variable assignment statement

		return lexAssignStmt

	} else if r == ';' {

		l.next()
		l.emit(ast.T_SEMICOLON)
		return lexStart

	} else if r == ',' {

		l.next()
		l.emit(ast.T_COMMA)
		return lexStart

	} else if r == '-' {

		// lex the slash prefix property name
		return lexProperty

	} else if r == '#' && r2 != '{' {

		return lexSelectors

	} else if unicode.IsLetter(r) || (r == '#') { // it might be -vendor- property or a property name or a selector

		// detect selector syntax
		l.remember()

		isProperty := false

		r = l.next()
		for r != EOF {
			// skip interpolation
			if r == '#' {
				if l.peek() == '{' {
					// find the matching brace
					r = l.next()
					for r != '}' {
						r = l.next()
					}
				}

			} else if r == ':' { // pseudo selector -> letters following ':', if there is a space after the ':' then it's a property value.

				if unicode.IsSpace(l.peek()) {
					isProperty = true
					break
				}

			} else if r == ';' {
				break
			} else if r == '}' {
				isProperty = true
				break
			} else if r == '{' {
				break
			} else if r == EOF {
				panic("unexpected EOF")
			}
			r = l.next()
		}

		l.rollback()

		if isProperty {
			return lexProperty
		} else {
			return lexSelectors
		}

	} else if r == '[' || r == '*' || r == '>' || r == '&' || r == '#' || r == '.' || r == '+' || r == ':' {

		return lexSelectors

	} else if r == EOF {

		return nil

	} else {

		l.errorf("Unexpected token: '%c'", r)

	}
	return nil
}

func lexStart(l *Lexer) stateFn {
	return lexStmt
}
