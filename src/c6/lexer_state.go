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
		// string start
		for {
			r = l.next()
			if r == '"' {
				l.emit(T_QQ_STRING)
				return lexStart
			} else if r == '\\' {
				// skip the escape character
				continue
			} else if r == EOF {
				panic("Expecting end of string")
			}
		}
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

func lexIdentifierSelector(l *Lexer) stateFn {
	var r = l.next()
	if r != '#' {
		l.error("Expecting '#' for lexing identifier, Got '%s'", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.error("An identifier should start with at least a letter, Got '%s'", r)
	}

	r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_ID_SELECTOR)

	// for selector like "#myID.foo"
	if r == '.' {
		l.emit(T_AND_SELECTOR)
	}
	return lexStatement
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

func lexAttributeSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '[' {
		l.emit(T_ATTRIBUTE_START)
		r = l.next()
		if !unicode.IsLetter(r) {
			l.error("Unexpected token for attribute name. Got '%s'", r)
		}
		for unicode.IsLetter(r) {
			r = l.next()
		}
		l.backup()
		l.emit(T_ATTRIBUTE_NAME)

		r = l.peek() // peek here again to avoid bugs
		if r == '=' {
			l.next()
			l.emit(T_EQUAL)

			r = l.peek()
			if r == '"' {
				lexString(l)
			} else {
				lexUnquoteString(l)
			}
		} else if l.match("~=") {
			l.emit(T_CONTAINS)
			r = l.peek()
			if r == '"' {
				lexString(l)
			} else {
				lexUnquoteString(l)
			}
		}

		r = l.peek()
		if r == ']' {
			l.next()
			l.emit(T_ATTRIBUTE_END)
			return lexStatement
		}

	}
	l.error("Unexpected token for attribute selector. Got '%s'", r)
	return nil
}

func lexClassSelector(l *Lexer) stateFn {
	var r = l.peek()
	if r == '.' {
		l.next()
		r = l.next()

		if !unicode.IsLetter(r) {
			panic("Expecting letter")
			return nil
		}

		for unicode.IsLetter(r) || r == '-' {
			r = l.next()
		}
		l.backup()
		l.emit(T_CLASS_SELECTOR)

		// there is a class name selector after this one.
		if r == '.' {
			l.emit(T_AND_SELECTOR)
		}
		return lexStart
	}
	return nil
}

func lexParentSelector(l *Lexer) stateFn {
	var r = l.peek()
	if r == '&' {
		l.next()
		l.emit(T_PARENT_SELECTOR)
		return lexStatement
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

func lexChildSelector(l *Lexer) stateFn {
	var r = l.peek()
	if r == '>' {
		l.emit(T_CHILD_SELECTOR)
		return lexSelector
	}
	l.error("Unexpected token '%s' for child selector.", r)
	return nil
}

func lexPseudoSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == ':' {
		r = l.next()

		if !unicode.IsLetter(r) {
			l.error("charater '%s' is not allowed in pseudo selector", r)
		}
		for unicode.IsLetter(r) || r == '-' {
			r = l.next()
		}
		l.backup()
		l.emit(T_PSEUDO_SELECTOR)

		if r == '(' {
			l.next()
			lexLang(l)
			r = l.next()
			if r != ')' {
				l.error("Unexpected token '%s' for pseudo lang selector", r)
			}
		}
		return lexStatement
	}
	l.error("Unexpected token '%s' for pseudo selector.", r)
	return nil
}

func lexUniversalSelector(l *Lexer) stateFn {
	var r = l.peek()
	if r == '*' {
		l.next()
		l.emit(T_UNIVERSAL_SELECTOR)
		return lexStatement
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

// Dispath selector lexing method
func lexSelector(l *Lexer) stateFn {
	l.ignoreSpaces()

	var r = l.peek()
	if unicode.IsLetter(r) {
		return lexTagNameSelector
	} else if r == '[' {
		return lexAttributeSelector
	} else if r == '.' {
		return lexClassSelector
	} else if r == '#' {
		return lexIdentifierSelector
	} else if r == '>' {
		return lexChildSelector
	} else if r == ':' {
		return lexPseudoSelector
	} else if r == '&' {
		return lexParentSelector
	} else if r == '*' {
		return lexUniversalSelector
	} else if r == ',' {

		l.next()
		l.emit(T_COMMA)

		// lex next selector
		return lexSelector

	} else if r == '+' {
		l.next()
		l.emit(T_ADJACENT_SELECTOR)
		return lexSelector
	} else if r == '{' {
		return nil
	} else {
		l.error("Unexpected token '%s' for selector.", r)
	}
	return nil
}

func lexTagNameSelector(l *Lexer) stateFn {
	var r = l.peek()
	if !unicode.IsLetter(r) {
		return lexStart
	}
	for unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_TAGNAME_SELECTOR)

	// predicate and inject the and selector for class name, identifier after the tagName
	if r == '.' || r == '#' || r == '[' || r == ':' {
		l.emit(T_AND_SELECTOR)
	}
	switch r {
	case ':':
		return lexPseudoSelector
	case '[':
		return lexAttributeSelector
	case '#':
		return lexIdentifierSelector
	case '.':
		return lexClassSelector
	}
	return lexStatement
}

func lexLang(l *Lexer) stateFn {
	/*
		html:lang(fr-ca) { quotes: '« ' ' »' }
		html:lang(de) { quotes: '»' '«' '\2039' '\203A' }
		:lang(fr) > Q { quotes: '« ' ' »' }
		:lang(de) > Q { quotes: '»' '«' '\2039' '\203A' }
	*/
	// [a-z]{2} - [a-z]{2}
	// [a-z]{2}
	var r = l.next()
	if !unicode.IsLetter(r) {
		l.error("Unexpected language token. Got '%s'", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.error("Unexpected language token. Got '%s'", r)
	}

	r = l.peek()
	if r == '-' {
		l.next() // skip '-'
		r = l.next()
		if !unicode.IsLetter(r) {
			l.error("Unexpected language token. Got '%s'", r)
		}
		r = l.next()
		if !unicode.IsLetter(r) {
			l.error("Unexpected language token. Got '%s'", r)
		}
	}
	l.emit(T_LANG_CODE)
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
	lexVariable(l)
	lexColon(l)
	fn := lexPropertyValue(l)
	return l.dispatchFn(fn)
}

func lexVariable(l *Lexer) stateFn {
	var r = l.next()
	if r != '$' {
		l.error("Unexpected token %s for lexVariable", r)
	}
	r = l.next()
	if !unicode.IsLetter(r) {
		l.error("The first character of a variable must be letter. Got '%s'", r)
	}

	r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_VARIABLE)
	return lexStatement
}

func lexExpansion(l *Lexer) stateFn {
	l.remember()
	var r rune = l.next()
	if r == '#' {
		r = l.next()
		if r == '{' {
			r = l.next()
			for r != '}' {
				r = l.next()
			}
			return nil
		}
	}
	l.rollback()
	return nil
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
	var r = l.peek()
	var r2 = l.peekMore(2)

	if r == 'p' && r2 == 'x' {
		l.advance(2)
		l.emit(T_UNIT_PX)
	} else if r == 'p' && r2 == 't' {
		l.advance(2)
		l.emit(T_UNIT_PT)
	} else if r == 'e' && r2 == 'm' {
		l.advance(2)
		l.emit(T_UNIT_EM)
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
func lexConstant(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_CONSTANT)
	return lexPropertyValue
}

func lexPropertyValue(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.peek()

	if r == '#' && l.peekMore(2) == '{' {
		return lexExpansion
	} else if r == '#' {
		return lexHexColor
	} else if unicode.IsDigit(r) {
		return lexNumber
	} else if unicode.IsLetter(r) {
		return lexConstant
	} else if r == '$' {
		return lexVariable
	} else if r == ' ' {
		l.next()
		l.ignore()
		return lexPropertyValue
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStatement
	} else if r == EOF {
		l.error("Unexpected end of file", r)
	} else {
		l.error("can't lex rune for property value: '%s'", r)
	}
	return nil
}

func lexColon(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r = l.next()
	if r == ':' {
		l.emit(T_COLON)
	} else {
		l.error("Expecting ':' token, Got '%s'", r)
	}
	return nil
}

func lexPropertyName(l *Lexer) stateFn {
	var r rune = l.next()
	for r == '-' || unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	l.emit(T_PROPERTY_NAME)
	lexColon(l)
	return lexPropertyValue
}

func lexStatement(l *Lexer) stateFn {
	l.ignoreSpaces()
	var r rune = l.peek()

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
	} else if r == '[' || r == '*' || r == '>' || r == '&' || r == '#' || r == '.' || r == '+' {

		return lexSelector
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
			} else if r == ';' {
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

			fmt.Println(string(r))
		}
	end_guess:

		// it's a selector, so we end with a brace '{'
		l.rollback()
		if isSelector {
			return lexSelector
		} else {
			return lexPropertyName
		}
	} else if r == ' ' {
		return lexSpaces
	} else if r == '"' || r == '\'' {
		return lexString
	} else if r == EOF {
		return nil
	} else {
		l.error("Can't lex rune: '%s'", r)
	}
	return nil
}

func lexStart(l *Lexer) stateFn {
	return lexStatement
}
