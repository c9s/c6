package lexer

import "unicode"
import _ "fmt"
import "github.com/c9s/c6/ast"

func IsInterpolationStartToken(r rune, r2 rune) bool {
	return r == '#' && r2 == '{'
}

// does not test ' '
func IsCombinatorToken(r rune) bool {
	return r == '>' || r == '+' || r == ',' || r == '~'
}

func IsSelector(t ast.TokenType) bool {
	return t == ast.T_CLASS_SELECTOR ||
		t == ast.T_ID_SELECTOR ||
		t == ast.T_TYPE_SELECTOR ||
		t == ast.T_UNIVERSAL_SELECTOR ||
		t == ast.T_PARENT_SELECTOR || // SASS parent selector
		t == ast.T_PSEUDO_SELECTOR || // :hover, :visited , ...
		t == ast.T_FUNCTIONAL_PSEUDO
}

/**
Pass peek() rune to check if it's a selector stop token
*/
func IsSelectorStopToken(r rune) bool {
	// pseudo, class, attribute, id, child, universal, adjacent
	return r == ':' ||
		r == '.' ||
		r == '[' ||
		r == '#' ||
		r == '&' ||
		r == '>' ||
		r == '*' ||
		r == '+' ||
		r == ','
}

func isDescendantCombinatorSeparator(r rune) bool {
	return r == ' '
}

func lexAttributeSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '[' {
		l.emit(ast.T_BRACKET_OPEN)

		var foundInterpolation = false

		r = l.next()
		if !unicode.IsLetter(r) && !IsInterpolationStartToken(r, l.peek()) {
			l.errorf("Unexpected token for attribute name. Got '%c'", r)
		}
		for {
			if IsInterpolationStartToken(r, l.peek()) {
				l.backup()
				lexInterpolation(l, false)
				foundInterpolation = true
			} else if !unicode.IsLetter(r) && r != '-' && r != '_' {
				break
			}
			r = l.next()
		}
		l.backup()

		token := l.createToken(ast.T_ATTRIBUTE_NAME)
		token.ContainsInterpolation = foundInterpolation
		l.emitToken(token)

		r = l.peek() // peek here again to avoid bugs

		var attrOp = false

		if r == '=' {

			l.next()
			l.emit(ast.T_ATTR_EQUAL)
			attrOp = true
		} else if l.match("~=") {

			l.emit(ast.T_INCLUDE_MATCH)
			attrOp = true
		} else if l.match("|=") {

			l.emit(ast.T_DASH_MATCH)
			attrOp = true
		} else if l.match("$=") {

			l.emit(ast.T_SUFFIX_MATCH)
			attrOp = true

		} else if l.match("*=") {
			l.emit(ast.T_SUBSTRING_MATCH)
			attrOp = true

		} else if l.match("^=") {
			l.emit(ast.T_PREFIX_MATCH)
			attrOp = true
		}

		if attrOp {
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
			l.emit(ast.T_BRACKET_CLOSE)
			return lexStmt
		}

	}
	l.errorf("Unexpected token for attribute selector. Got '%c'", r)
	return nil
}

func lexClassSelector(l *Lexer) stateFn {
	l.accept(".")

	var r = l.next()
	if !unicode.IsLetter(r) {
		l.errorf("Expecting letter for class selector. got '%c'", r)
		return nil
	}

	// skip valid class name characters
	for unicode.IsLetter(r) || r == '-' || r == '_' {
		r = l.next()
	}
	l.backup()
	l.emit(ast.T_CLASS_SELECTOR)
	return lexSelectors
}

func lexPseudoSelector(l *Lexer) stateFn {
	var foundInterpolation = false

	// the first ':'
	var r = l.next()

	// support CSS3 syntax for `::before` and `::after`
	// @see https://developer.mozilla.org/en-US/docs/Web/CSS/::before
	l.accept(":")

	r = l.next()
	if !unicode.IsLetter(r) && !(r == '#' && l.peek() == '{') {
		l.errorf("charater '%c' is not allowed in pseudo selector", r)
	}
	for r != EOF && (unicode.IsLetter(r) || r == '-' || r == '#') {
		if IsInterpolationStartToken(r, l.peek()) {
			l.backup()
			lexInterpolation(l, false)
			foundInterpolation = true
		}
		r = l.next()
	}
	l.backup()

	if foundInterpolation {
		l.emit(ast.T_INTERPOLATION_SELECTOR)
	} else {

		r = l.peek()
		if r == '(' {
			l.emit(ast.T_FUNCTIONAL_PSEUDO)
			l.next()
			l.emit(ast.T_PAREN_OPEN)

			r = l.peek()
			for r != ')' && r != EOF {
				if fn := lexExpr(l); fn == nil {
					break
				}
				r = l.peek()
			}

			l.expect(")")
			l.emit(ast.T_PAREN_CLOSE)

		} else {
			l.emit(ast.T_PSEUDO_SELECTOR)
		}
	}
	return lexSelectors
}

func lexUniversalSelector(l *Lexer) stateFn {
	l.expect("*")
	l.emit(ast.T_UNIVERSAL_SELECTOR)
	return lexSimpleSelector
}

func lexSimpleSelector(l *Lexer) stateFn {

	var r = l.peek()

	if r == '.' {

		return lexClassSelector

	} else if r == '[' {

		return lexAttributeSelector

	} else if r == ':' {

		return lexPseudoSelector

	} else if r == '#' && l.peekBy(2) != '{' {

		return lexIdSelector

	} else if r == '&' {

		l.next()
		l.emit(ast.T_PARENT_SELECTOR)
		return lexSelectors

	} else if r == '*' {

		return lexUniversalSelector

	} else if unicode.IsLetter(r) {

		return lexTypeSelector

	}

	return lexSelectors
}

// Dispath selector lexing method
func lexSelectors(l *Lexer) stateFn {
	var r rune

	lexComment(l, false)

	// space between selector means descendant selector
	if tok := l.lastToken(); tok != nil && IsSelector(tok.Type) {
		var foundSpace = false
		var r = l.next()
		for unicode.IsSpace(r) || r == '/' {
			if unicode.IsSpace(r) {
				foundSpace = true
			}
			lexComment(l, false)
			r = l.next()
		}
		l.backup()
		if r == EOF {
			return nil
		}
		if foundSpace && r != ',' && r != '{' && !IsCombinatorToken(r) {
			l.emit(ast.T_DESCENDANT_COMBINATOR)
		} else {
			l.ignore()
		}
	}

	lexComment(l, false)

	// re-peek again
	r = l.peek()

	// lex the first selector
	if unicode.IsLetter(r) {

		return lexTypeSelector

	} else if r == '[' {

		return lexAttributeSelector

	} else if r == '.' {

		return lexClassSelector

	} else if r == ':' {

		return lexPseudoSelector

	} else if r == '&' {

		l.expect("&")
		l.emit(ast.T_PARENT_SELECTOR)
		return lexSelectors

	} else if r == '*' {

		return lexUniversalSelector

	} else if r == '#' {
		// for selector syntax like:
		//    '#{  }  {  }'
		//    '#{ a }foo#{ b } {  }'
		//    '#{  }.something {  }'
		//    '#{  } .something {  }'
		//    '#{  }#myId {  }'
		if IsInterpolationStartToken(r, l.peekBy(2)) {
			if tok := l.lastToken(); tok != nil && IsSelector(tok.Type) {
				l.emit(ast.T_LITERAL_CONCAT)
			}

			lexInterpolation(l, false)
			// end of interpolation

			// find stop point of a selector.
			r = l.next()
			for {
				if IsInterpolationStartToken(r, l.peek()) {
					l.backup()
					lexInterpolation(l, false)
				} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' || IsSelectorStopToken(r) || isDescendantCombinatorSeparator(r) {
					break
				}
				r = l.next()
			}
			l.backup()

			// the suffix of the selector.
			var token = l.createToken(ast.T_INTERPOLATION_SELECTOR)
			token.ContainsInterpolation = true
			l.emitToken(token)
			return lexSelectors
		} else {
			return lexIdSelector
		}

	} else if r == '>' {

		l.next()
		l.emit(ast.T_CHILD_COMBINATOR)
		return lexSelectors

	} else if r == '~' {

		l.next()
		l.emit(ast.T_GENERAL_SIBLING_COMBINATOR)
		return lexSelectors

	} else if r == ',' {

		l.next()
		l.emit(ast.T_COMMA)

		// lex next selector
		return lexSelectors

	} else if r == '+' {

		l.next()
		l.emit(ast.T_ADJACENT_SIBLING_COMBINATOR)
		return lexSelectors

	} else if unicode.IsSpace(r) {

		l.next()
		for unicode.IsSpace(r) {
			r = l.next()
		}
		l.backup()
		l.ignore()

		return lexSelectors

	} else if r == '{' || r == ';' {

		return lexStmt

	} else {
		l.errorf("Unexpected token '%c' for lexing selector.", r)
	}
	return nil
}

func lexTypeSelector(l *Lexer) stateFn {
	var r = l.next()
	if !unicode.IsLetter(r) && !IsInterpolationStartToken(r, l.peekBy(2)) {
		l.errorf("Expecting letter token for tag name selector. got %c", r)
	}

	var foundInterpolation = false
	r = l.next()
	for {
		if IsInterpolationStartToken(r, l.peek()) {
			l.backup()
			lexInterpolation(l, false)
			foundInterpolation = true
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}
		r = l.next()
	}
	l.backup()

	if foundInterpolation {
		l.emit(ast.T_INTERPOLATION_SELECTOR)
	} else {
		l.emit(ast.T_TYPE_SELECTOR)
	}

	return lexSimpleSelector
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
		l.errorf("Unexpected language token. Got '%c'", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.errorf("Unexpected language token. Got '%c'", r)
	}

	r = l.peek()
	if r == '-' {
		l.next() // skip '-'
		r = l.next()
		if !unicode.IsLetter(r) {
			l.errorf("Unexpected language token. Got '%c'", r)
		}
		r = l.next()
		if !unicode.IsLetter(r) {
			l.errorf("Unexpected language token. Got '%c'", r)
		}
	}
	l.emit(ast.T_LANG_CODE)
	return nil
}

func lexIdSelector(l *Lexer) stateFn {
	var foundInterpolation = false
	var r = l.next()
	r = l.next()
	if !unicode.IsLetter(r) && r != '#' && l.peek() != '{' {
		l.errorf("An identifier should start with at least a letter, Got '%c'", r)
	}
	for {
		if IsInterpolationStartToken(r, l.peek()) {
			l.backup()
			lexInterpolation(l, false)
			foundInterpolation = true
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}
		r = l.next()
	}
	l.backup()

	r = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()

	if foundInterpolation {
		l.emit(ast.T_INTERPOLATION_SELECTOR)
	} else {
		l.emit(ast.T_ID_SELECTOR)
	}
	return lexSelectors
}
