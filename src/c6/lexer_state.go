package c6

import (
	// "fmt"
	"unicode"
)

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
		l.next()
		for {
			r = l.next()
			if r == '"' {
				l.next()
				l.emit(T_QQ_STRING)
				return lexStart
			} else if r == '\\' {
				// skip the escape character
				l.next()
			} else if r == eof {
				panic("Expecting end of string")
			}
		}
		return lexStart

	} else if r == '\'' {

		return lexStart

	}
	l.backup()
	return nil
}

func lexAtRule(l *Lexer) stateFn {
	t := l.next()
	if t == '@' {
		if l.match("import") {
			l.emit(T_AT_IMPORT)
			return lexSpaces
		}
	}
	l.backup()
	return nil
}

func lexStatement(l *Lexer) stateFn {
	var t = l.peek()
	if t == '/' {
		return lexComment
	} else if t == '@' {

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

func lexStart(l *Lexer) stateFn {
	l.ignoreSpaces()
	// var r rune = l.peek()
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
