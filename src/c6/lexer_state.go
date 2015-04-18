package c6

import "unicode"
import "strings"
import "fmt"

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

func (l *Lexer) emitIfKeywordMatches() bool {
	l.remember()
	for keyword, typ := range LexKeywords {
		var match bool = true
	next_keyword:
		for _, sc := range keyword {
			c := l.next()
			if c == eof {
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
			if c == '\n' || c == eof || c == ' ' || c == '\t' || unicode.IsSymbol(c) {
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
			} else if r == eof {
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
			} else if r == eof {
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
	var r rune = l.peek()
	if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStart
	}
	return nil
}

func lexExpansionStart(l *Lexer) stateFn {

	return nil
}

func lexHexColor(l *Lexer) stateFn {
	l.ignoreSpaces()
	l.remember()

	var r rune = l.next()
	if r == '#' {
		l.next()
		for l.accept("0123456789abcdefABCDEF") {
			l.next()
		}
		if (l.Offset-l.Start) != 4 && (l.Offset-l.Start) != 7 {
			panic("Invalid hex color length")
		}
		l.emit(T_HEX_COLOR)
		return lexPropertyValue
	}

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
	if r == '#' && l.peekMore(2) == '{' {
		return lexExpansionStart
	} else if r == '#' {
		return lexHexColor
	} else if unicode.IsDigit(r) {
		return lexDigits
	} else if unicode.IsLetter(r) {
		return lexConstant
	} else if r == ' ' {
		l.next()
		l.ignore()
		return lexPropertyValue
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStatementOrProperty
	} else {
		panic(fmt.Sprintf("can't lex rune: %+v", string(r)))
	}
}

func lexPropertyName(l *Lexer) stateFn {
	var r rune = l.next()
	for r == '-' || unicode.IsLetter(r) {
		r = l.next()
	}

	// correct stop
	if r != ':' {
		panic("invalid property name")
	}
	l.backup()
	l.emit(T_PROPERTY_NAME)
	l.next() // skip ":"
	return lexPropertyValue
}

// there is a possibility that a statement start with "width:", "input{", "input {"
func lexStatementOrProperty(l *Lexer) stateFn {
	l.ignoreSpaces()

	// var r rune = l.peek()
	var str = l.lookaheadTil(" ,:;{")
	// fmt.Println("lexStatementOrProperty", str)

	// looks like property
	if strings.HasSuffix(str, ":") {
		l.Offset += len(str)
		l.emit(T_PROPERTY_NAME)
		return lexPropertyValue
	} else {
		return lexStatement
	}
	return nil

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
		return lexStatementOrProperty
	} else if r == '}' {
		l.next()
		l.emit(T_BRACE_END)
		return lexStatement
	} else if r == ';' {
		l.next()
		l.emit(T_SEMICOLON)
		return lexStart
	} else if r == '@' {
		return lexAtRule
	} else if unicode.IsLetter(r) { // it maybe -vendor- property or a property name
		return lexTagName
	} else if r == '.' {
		return lexClassName
	} else if r == ' ' {
		return lexSpaces
	} else if r == '"' || r == '\'' {
		return lexString
	} else if r == eof {
		return nil
	} else {
		panic(fmt.Sprintf("can't lex rune: %+v", string(r)))
	}
}

func lexStart(l *Lexer) stateFn {
	return lexStatement
	/*
		if unicode.IsDigit(c) {
			return lexNumber
		} else if l.emitIfMatch("->", T_FUNCTION_GLYPH) {
			return lexStart
		} else if l.emitIfMatch("::", T_FUNCTION_PROTOTYPE) {
			return lexStart
		} else if l.consumeIfMatch("//") {
			return lexOnelineComment
		} else if l.consumeIfMatch("/*") {
			return lexComment
		} else if c == '-' {
			l.next()
			l.emit(TokenType(c))
			return lexStart
		} else if c == ' ' || c == '\t' {
			// return lexSpaces
			return lexIgnoreSpaces
		} else if c == '\n' || c == '\r' {
			// if there is a new line, check the next line is indent or outdent,
			// if there is no spaces/indent in the next line, then it should be outdent.
			l.Line++
			c = l.next() // take the the line break char

			// skip multiple newline at one time
			// sometimes we wrote:
			// --->a = 3$
			// $
			// --->b = 10$
			// and this should be in the same block.
			for c == '\n' {
				c = l.next()
			}
			l.backup()

			// c = l.peek()
			if c == eof {
				return lexStart
			}

			// reset space info
			l.lastSpace = l.space
			l.space = 0

			// calculate spaces
			c = l.next()
			for c == ' ' || c == '\t' {
				l.space++
				c = l.next()
			}
			l.backup()
			if l.space == l.lastSpace {
				l.emit(T_NEWLINE)
			} else if l.space < l.lastSpace {
				l.emit(T_OUTDENT)
				l.emit(T_NEWLINE) // means end of statement
				l.IndentLevel--
			} else if l.space > l.lastSpace {
				l.IndentLevel++
				l.emit(T_INDENT)
			}
			return lexStart
		} else if l.emitIfMatch("&&", T_BOOLEAN_AND) || l.emitIfMatch("||", T_BOOLEAN_OR) {
			return lexStart
		} else if l.emitIfMatch("==", T_EQUAL) || l.emitIfMatch(">=", T_GT_EQUAL) || l.emitIfMatch("<=", T_LT_EQUAL) {
			return lexStart
		} else if l.emitIfKeywordMatches() {
			return lexStart
		} else if l.accept("+-*|&[]{}()<>,=@") {
			l.emit(TokenType(c))
			return lexStart
		} else if l.lastTokenType == T_NUMBER && l.emitIfMatch("..", T_RANGE_OPERATOR) {
			return lexStart
		} else if c == '"' || c == '\'' {
			return lexString
		} else if unicode.IsLetter(c) {
			return lexIdentifier
		} else if c == eof {
			if l.IndentLevel > 0 {
				for i := 0; i < l.IndentLevel; i++ {
					l.emit(T_OUTDENT)
					l.emit(T_NEWLINE)
				}
			}
			return nil
		} else {
			panic(fmt.Errorf("unknown token %c\n", c))
			return nil
		}
	*/
	return nil
}
