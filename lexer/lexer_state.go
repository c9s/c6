package lexer

import (
	"fmt"
	"github.com/c9s/c6/ast"
	"unicode"
)

type stateFn func(*Lexer) stateFn

const LETTERS = "zxcvbnmasdfghjklqwertyuiop"
const DIGITS = "1234567890"

func (l *Lexer) errorf(msg string, r rune) {
	var err = fmt.Errorf(msg, string(r))
	panic(err)
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
	l.ignoreComment()
	l.ignoreSpaces()

	if l.accept(";") {
		l.emit(ast.T_SEMICOLON)
	} else if l.accept("}") {
		l.emit(ast.T_BRACE_CLOSE)
	}
	return lexStart
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
	return lexStart
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
		return lexStart
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
