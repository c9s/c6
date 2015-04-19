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

func lexClassName(l *Lexer) stateFn {
	var t = l.peek()
	if t == '.' {
		l.next()
		t = l.next()

		if !unicode.IsLetter(t) {
			panic("Expecting letter")
			return nil
		}

		for unicode.IsLetter(t) || t == '-' {
			t = l.next()
		}
		l.backup()
		l.emit(T_CLASS_SELECTOR)
		return lexStart
	}
	return nil
}

func lexTagName(l *Lexer) stateFn {
	var t = l.peek()
	if !unicode.IsLetter(t) {
		return lexStart
	}
	for unicode.IsLetter(t) {
		t = l.next()
	}
	l.backup()
	l.emit(T_TAGNAME_SELECTOR)
	return lexStart
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

func lexVariable(l *Lexer) stateFn {
	var r = l.next()
	if r != '$' {
		l.error("Unexpected token %s for lexVariable", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		fmt.Printf("The first character of a variable must be letter.")
	}
	r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		l.next()
	}
	l.backup()
	l.emit(T_VARIABLE)
	lexColon(l)
	return lexPropertyValue
}

func lexExpansionStart(l *Lexer) stateFn {

	// TODO
	return nil
}

func lexHexColor(l *Lexer) stateFn {
	l.ignoreSpaces()
	l.remember()

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
	l.error("expecting hex color", r)
	return nil
}

func lexDigits(l *Lexer) stateFn {
	var r = l.next()
	for unicode.IsDigit(r) {
		l.next()
	}
	l.backup()
	l.emit(T_DIGITS)
	return nil
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
	fmt.Printf("lexPropertyValue: %s\n", string(r))
	if r == '#' && l.peekMore(2) == '{' {
		return lexExpansionStart
	} else if r == '#' {
		return lexHexColor
	} else if unicode.IsDigit(r) {
		return lexDigits
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
	} else {
		panic(fmt.Sprintf("can't lex rune: %+v", string(r)))
	}
}

func lexColon(l *Lexer) stateFn {
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
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStart
	} else if r == '$' {
		return lexVariable
	} else if r == '@' {
		return lexAtRule
	} else if r == '-' || unicode.IsLetter(r) { // it maybe -vendor- property or a property name
		l.remember()

		// if it starts with a letter, it's possible to have two kinds of syntax here:
		//   a { }
		//   a{}
		//   a:hover {  }
		//   color: red;
		//   background-color: ...
		//   -webkit-transition: ...

		// lex the property name (or tag name)
		for l.accept(LETTERS + "-") {
		}
		// ignore spaces and colon
		l.accept(": ")

		var r = l.next()
		for unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			r = l.next()
		}

		// it's a selector, so we end with a brace '{'
		l.rollback()
		if r == '{' {
			return lexTagName
		} else if r == ';' {
			return lexPropertyName
		} else {
			return lexPropertyName
		}
	} else if r == '.' {
		return lexClassName
	} else if r == ' ' {
		return lexSpaces
	} else if r == '"' || r == '\'' {
		return lexString
	} else if r == EOF {
		return nil
	} else {
		l.error("can't lex rune: %+v", r)
	}
	return nil
}

func lexStart(l *Lexer) stateFn {
	return lexStatement
}
