package c6

import "unicode"

// import "strings"
import "fmt"
import "errors"
import "c6/ast"

type stateFn func(*Lexer) stateFn

const LETTERS = "zxcvbnmasdfghjklqwertyuiop"
const DIGITS = "1234567890"

var LexKeywords = map[string]int{}

func (l *Lexer) error(msg string, r rune) {
	var err = errors.New(fmt.Sprintf(msg, string(r)))
	panic(err)
}

func (l *Lexer) emitIfKeywordMatches() bool {
	l.remember()
	for keyword, typ := range LexKeywords {
		var match bool = true
	next_keyword:
		for _, sc := range keyword {
			c := l.next()
			if c == EOF {
				match = false
				break next_keyword
			}
			if sc != c {
				match = false
				break next_keyword
			}
		}
		if match {
			c := l.next()
			if c == '\n' || c == EOF || c == ' ' || c == '\t' || unicode.IsSymbol(c) {
				l.backup()
				l.emit(ast.TokenType(typ))
				return true
			}
		}
		l.rollback()
	}
	return false
}

func lexComment(l *Lexer) stateFn {
	var r = l.next()

	if r == '/' {
		r = l.next()
		if r == '/' {
			for {
				r = l.next()
				if r == '\n' || r == '\r' {
					l.emit(ast.T_COMMENT_LINE)
					return lexStart
				}
			}
		} else if r == '*' {

			for {
				r = l.next()
				if r == '*' && l.peek() == '/' {
					l.next()
					l.emit(ast.T_COMMENT_BLOCK)
					return lexStart
				}
			}
		} else {
			l.backup()
		}
	}
	l.backup()
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
				continue
			} else if isInterpolationStartToken(r, l.peek()) {
				l.backup()
				lexInterpolation(l, false)
				containsInterpolation = true
			} else if r == EOF {
				panic("Expecting end of string")
			}
			r = l.next()
		}
		l.backup()
		return lexStart

	} else if r == '\'' {

		l.next()
		l.ignore()
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
		return lexStart
	}
	l.backup()
	return nil
}

func lexUrl(l *Lexer) {
	if l.match("url") {
		l.emit(ast.T_IDENT)
		l.match("(")
		l.emit(ast.T_PAREN_START)

		var q = l.peek()
		if q == '"' || q == '\'' {
			lexString(l)
		} else {
			lexUnquoteStringStopAt(l, ')')
		}
		l.match(")")
		l.emit(ast.T_PAREN_END)

	} else {
		var r = l.peek()
		if r == '"' || r == '\'' {
			lexString(l)
		} else {
			l.error("Unexpected token for @import rule. Got %s", r)
		}
	}
}

/*
func lexMediaQuery(l *Lexer) stateFn {
	if !unicode.IsLetter(l.peek()) {
		return nil
	}

	var r = l.next()
	for {
		r = l.next()
	}
	l.backup()
	for unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.ignoreSpaces()

	if l.peek() == ',' {
		l.next()
		l.emit(T_COMMA)
	}
}
*/

/*

Currently the @import rule only supports '@import url(...) media;

@see https://developer.mozilla.org/en-US/docs/Web/CSS/@import for more @import syntax support
*/
func lexAtRule(l *Lexer) stateFn {
	t := l.peek()

	if t == '@' {
		l.next()
		if l.match("import") {
			l.emit(ast.T_IMPORT)
			l.ignoreSpaces()

			lexUrl(l)
			l.ignoreSpaces()

			// looks like a media list
			for unicode.IsLetter(l.peek()) {
				l.next()
			}
			if l.precedeStartOffset() {
				l.emit(ast.T_MEDIA)
			}
			return lexStatement
		} else if l.match("charset") {
			l.emit(ast.T_CHARSET)
			l.ignoreSpaces()
			return lexStatement
		} else {
			panic("Unknown at-rule directive")
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
	return lexStart
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

func lexSemiColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.next()
	if r == ';' {
		l.emit(ast.T_SEMICOLON)
		return lexStatement
	}
	l.backup()
	return nil
}

func lexVariableAssignment(l *Lexer) stateFn {
	lexVariableName(l)
	lexColon(l)
	var r = l.peek()
	for r != ';' && r != '}' {
		lexExpression(l)
		r = l.peek()
	}
	r = l.next()
	if r == ';' {
		l.emit(ast.T_SEMICOLON)
	} else if r == '}' {
		l.emit(ast.T_BRACE_END)
	}
	return lexStatement
}

// $var-rgba(255,255,0)
func lexVariableName(l *Lexer) stateFn {
	var r = l.next()
	if r != '$' {
		l.error("Unexpected token %s for lexVariable", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.error("The first character of a variable name must be letter. Got '%s'", r)
	}

	r = l.next()
	for {
		if r == '-' {
			var r2 = l.peek()
			if unicode.IsLetter(r2) { // $a-b is a valid variable name.
				l.next()
			} else if unicode.IsDigit(r2) { // $a-3 should be $a '-' 3
				l.backup()
				l.emit(ast.T_VARIABLE)
				return lexExpression
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
			return lexStatement
			break
		} else if r == EOF || r == ' ' || r == ';' {
			break
		}
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_VARIABLE)
	return lexStatement
}

func lexHexColor(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.next()
	if r != '#' {
		l.error("Expecting hex color, got '%s'", r)
	}

	r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()

	var length = l.length() - 1
	if length != 3 && length != 6 {
		panic(fmt.Errorf("Invalid hex color length, expecting 3 or 6, got %d - %s", length, l.current()))
	}
	l.emit(ast.T_HEX_COLOR)
	return lexExpression
}

func lexNumberUnit(l *Lexer) stateFn {
	if l.match("px") {
		l.emit(ast.T_UNIT_PX)
	} else if l.match("pt") {
		l.emit(ast.T_UNIT_PT)
	} else if l.match("em") {
		l.emit(ast.T_UNIT_EM)
	} else if l.match("cm") {
		l.emit(ast.T_UNIT_CM)
	} else if l.match("mm") {
		l.emit(ast.T_UNIT_MM)
	} else if l.match("rem") {
		l.emit(ast.T_UNIT_REM)
	} else if l.match("deg") {
		l.emit(ast.T_UNIT_DEG)
	} else if l.match("%") {
		l.emit(ast.T_UNIT_PERCENT)
	} else if l.peek() == ';' {
		return lexStatement
	}
	return lexExpression
}

func lexNumber(l *Lexer) stateFn {
	var r = l.next()

	var floatPoint = false

	for unicode.IsDigit(r) {
		r = l.next()
		if r == '.' {
			floatPoint = true
			r = l.next()
			if !unicode.IsDigit(r) {
				l.error("Expecting at least one digit after the floating point, got '%s'", r)
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

func lexStatement(l *Lexer) stateFn {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r rune = l.peek()

	if r == EOF {
		return nil
	}

	if r == '@' {
		return lexAtRule
	} else if r == '(' {
		l.next()
		l.emit(ast.T_PAREN_START)
		return lexStart
	} else if r == ')' {
		l.next()
		l.emit(ast.T_PAREN_END)
		return lexStart
	} else if r == '{' {
		l.next()
		l.emit(ast.T_BRACE_START)
		return lexStatement
	} else if r == '}' {
		l.next()
		l.emit(ast.T_BRACE_END)
		return lexStatement
	} else if r == '/' && l.peekBy(2) == '/' {
		l.next()
		r = l.next()
		for r != '\r' && r != '\n' && r != EOF {
			l.next()
		}

	} else if r == '$' { // it's a variable assignment statement
		return lexVariableAssignment
	} else if r == '[' || r == '*' || r == '>' || r == '&' || r == '#' || r == '.' || r == '+' || r == ':' {
		return lexSelectors
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

	} else if unicode.IsLetter(r) { // it might be -vendor- property or a property name or a selector

		// detect selector syntax
		l.remember()

		isSelector := false

		r = l.next()
		for {
			// ignore interpolation
			if r == '#' && l.peekBy(2) == '{' {
				// find the matching brace
				r = l.next()
				for r != '}' {
					r = l.next()
				}
			} else if r == '{' {
				isSelector = true
				break
			} else if r == ';' {
				isSelector = false
				break
			} else if r == '}' {
				isSelector = false
				break
			} else if r == EOF {
				panic("unexpected EOF")
				break
			}
			r = l.next()
		}

		// it's a selector, so we end with a brace '{'
		l.rollback()
		if isSelector {
			return lexSelectors
		} else {
			return lexProperty
		}
	} else if r == '"' || r == '\'' {
		return lexString
	} else if r == EOF {
		return nil
	} else {
		l.error("Can't lex rune in lexStatement: '%s'", r)
	}
	return nil
}

func lexStart(l *Lexer) stateFn {
	return lexStatement
}
