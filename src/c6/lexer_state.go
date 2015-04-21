package c6

import "unicode"

// import "strings"
import "fmt"
import "errors"

type stateFn func(*Lexer) stateFn

const LETTERS = "zxcvbnmasdfghjklqwertyuiop"
const DIGITS = "1234567890"

var LexKeywords = map[string]int{
/*
	"if":      T_IF,
	"class":   T_CLASS,
	"for":     T_FOR,
	"foreach": T_FOREACH,
	"else":    T_ELSE,
	"elseif":  T_ELSEIF,
	"echo":    T_ECHO,
	"is":      T_IS,
	"return":  T_RETURN,
	"extends": T_EXTENDS,
	"does":    T_DOES,
	"new":     T_NEW,
	"clone":   T_CLONE,
*/
}

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
				l.emit(TokenType(typ))
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
					l.emit(T_COMMENT_LINE)
					return lexStart
				}
			}
		} else if r == '*' {

			for {
				r = l.next()
				if r == '*' && l.peek() == '/' {
					l.next()
					l.emit(T_COMMENT_BLOCK)
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

		// string start
		r = l.next()
		for {
			if r == '"' {
				token := l.createToken(T_QQ_STRING)
				token.ContainsInterpolation = containsInterpolation
				l.emitToken(token)
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
		for {
			r = l.next()
			if r == '\'' {
				l.next()
				l.emit(T_Q_STRING)
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

func lexAtRule(l *Lexer) stateFn {
	t := l.peek()
	// fmt.Printf("%c", t)
	if t == '@' {
		l.next()
		if l.match("import") {
			// fmt.Printf("match @import")
			l.emit(T_IMPORT)
			return lexStart
		} else if l.match("charset") {
			l.emit(T_CHARSET)
			return lexStart
		} else {
			panic("unknown at-rule directive")
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

func lexUnquoteString(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_UNQUOTE_STRING)
	return nil
}

func lexSemiColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.next()
	if r == ';' {
		l.emit(T_SEMICOLON)
		return lexStatement
	}
	l.backup()
	return nil
}

func lexVariableAssignment(l *Lexer) stateFn {
	lexVariableName(l)
	lexColon(l)
	return lexPropertyValue(l)
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
				l.emit(T_VARIABLE)
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
			l.emit(T_VARIABLE)
			return lexStatement
			break
		} else if r == EOF || r == ' ' || r == ';' {
			break
		}
		r = l.next()
	}
	l.backup()
	l.emit(T_VARIABLE)
	return lexStatement
}

func lexHexColor(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.next()
	if r == '#' {
		for l.accept("0123456789abcdefABCDEF") {
		}
		if (l.Offset-l.Start) != 4 && (l.Offset-l.Start) != 7 {
			panic("Invalid hex color length")
		}
		l.emit(T_HEX_COLOR)
		return lexPropertyValue
	}
	l.error("Expecting hex color, got '%s'", r)
	return nil
}

func lexNumberUnit(l *Lexer) stateFn {
	if l.match("px") {
		l.emit(T_UNIT_PX)
	} else if l.match("pt") {
		l.emit(T_UNIT_PT)
	} else if l.match("em") {
		l.emit(T_UNIT_EM)
	} else if l.match("deg") {
		l.emit(T_UNIT_DEG)
	} else if l.match("%") {
		l.emit(T_UNIT_PERCENT)
	} else if l.peek() == ';' {
		return lexStatement
	}
	return lexPropertyValue
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
		l.emit(T_FLOAT)
	} else {
		l.emit(T_INTEGER)
	}
	return lexNumberUnit
}

// lex for: `center`, `auto`, `top`, `none`
func lexConstantString(l *Lexer) stateFn {
	var r = l.next()

	// first char should be letter
	if !unicode.IsLetter(r) {
		l.error("Unexpected token for constant string. Got '%s'", r)
	}

	r = l.next()
	for unicode.IsLetter(r) || r == '-' {
		r = l.next()
	}
	l.backup()
	l.emit(T_CONSTANT)
	return lexPropertyValue
}

func lexStatement(l *Lexer) stateFn {
	// strip the leading spaces of a statement
	l.ignoreSpaces()

	var r rune = l.peek()

	if r == EOF {
		return nil
	}

	if r == '(' {
		l.next()
		l.emit(T_PAREN_START)
		return lexStart
	} else if r == ')' {
		l.next()
		l.emit(T_PAREN_END)
		return lexStart
	} else if r == '{' {
		l.next()
		l.emit(T_BRACE_START)
		return lexStatement
	} else if r == '}' {
		l.next()
		l.emit(T_BRACE_END)
		return lexStatement
	} else if r == '$' { // it's a variable assignment statement
		return lexVariableAssignment
	} else if r == '[' || r == '*' || r == '>' || r == '&' || r == '#' || r == '.' || r == '+' || r == ':' {
		return lexSelectors
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStart
	} else if r == ',' {
		l.next()
		l.emit(T_COMMA)
		return lexStart
	} else if r == '@' {
		return lexAtRule
	} else if r == '-' || unicode.IsLetter(r) { // it maybe -vendor- property or a property name

		// detect selector syntax
		l.remember()

		isSelector := false

		r = l.next()
		for {
			if r == EOF {
				break
			}
			if !unicode.IsLetter(r) && r != '-' {
				break
			}

			// ignore interpolation
			if r == '#' && l.peekMore(2) == '{' {
				r = l.next()
				// find the matching brace
				for r != '}' {
					r = l.next()
				}
			}

			// for unicode.IsLetter(r) || r == '-' || r == ' ' {
			r = l.next()
		}

		for r == ' ' {
			r = l.next()
		}
		if r == '{' {
			isSelector = true
			goto end_guess
		}

		if r != ':' {
			isSelector = true
			goto end_guess
		}

		r = l.next()
		// ignore space
		for r == ' ' {
			r = l.next()
		}
		for {
			if r == '{' {
				isSelector = true
				goto end_guess
			} else if r == '}' { // end of property
				isSelector = false
				goto end_guess
			} else if r == ';' {
				isSelector = false
				goto end_guess
			} else if r == EOF {
				isSelector = false
				goto end_guess
			} else if r == '#' && l.peekMore(2) == '{' {
				// skip expansion
				r = l.next()
				for r != '}' {
					r = l.next()
				}
				l.backup()
			}
			r = l.next()
		}
	end_guess:

		// it's a selector, so we end with a brace '{'
		l.rollback()
		if isSelector {
			return lexSelectors
		} else {
			return lexPropertyName
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
