package c6

import "unicode"

func isInterpolationStartToken(r rune, r2 rune) bool {
	return r == '#' && r2 == '{'
}

func isSelector(t TokenType) bool {
	return t == T_CLASS_SELECTOR ||
		t == T_ID_SELECTOR ||
		t == T_CLASS_SELECTOR ||
		t == T_TAGNAME_SELECTOR ||
		t == T_UNIVERSAL_SELECTOR ||
		t == T_PARENT_SELECTOR || // SASS parent selector
		t == T_PSEUDO_SELECTOR || // :hover, :visited , ...
		t == T_INTERPOLATION_SELECTOR // selector contains interpolation: '#{ ... }'
}

/**
Pass peek() rune to check if it's a selector stop token
*/
func isSelectorStopToken(r rune) bool {
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

func isDescendantSelectorSeparator(r rune) bool {
	return r == ' '
}

func lexAttributeSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '[' {
		l.emit(T_BRACKET_LEFT)
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

		var attrOp = false

		if r == '=' {
			l.next()
			l.emit(T_EQUAL)
			attrOp = true
		} else if l.match("~=") {
			l.emit(T_TILDE_EQUAL)
			attrOp = true
		} else if l.match("|=") {
			l.emit(T_PIPE_EQUAL)
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
			l.emit(T_BRACKET_RIGHT)
			return lexStatement
		}

	}
	l.error("Unexpected token for attribute selector. Got '%s'", r)
	return nil
}

func lexClassSelector(l *Lexer) stateFn {
	var r = l.next()
	if r != '.' {
		l.error("Unexpected token for class selector. got '%s'", r)
	}

	r = l.next()
	if !unicode.IsLetter(r) {
		l.error("Expecting letter for class selector. got '%s'", r)
		return nil
	}

	// skip valid class name characters
	for unicode.IsLetter(r) || r == '-' || r == '_' {
		r = l.next()
	}
	l.backup()
	l.emit(T_CLASS_SELECTOR)
	return lexSelectors
}

func lexParentSelector(l *Lexer) stateFn {
	var r = l.peek()
	if r == '&' {
		l.next()
		l.emit(T_PARENT_SELECTOR)

		if l.accept(".:[") {
			l.backup()
			l.emit(T_AND_SELECTOR)
		}
		return lexSelectors
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

func lexChildSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '>' {
		l.emit(T_GT)
		return lexSelectors
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
			l.ignore()
			lexLang(l)
			r = l.next()
			if r != ')' {
				l.error("Unexpected token '%s' for pseudo lang selector", r)
			}
			l.ignore()
		}
		return lexSelectors
	}
	l.error("Unexpected token '%s' for pseudo selector.", r)
	return nil
}

func lexUniversalSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '*' {
		l.emit(T_UNIVERSAL_SELECTOR)

		r = l.peek()
		if r == '.' {
			l.emit(T_AND_SELECTOR)
			return lexClassSelector
		} else if r == '[' {
			l.emit(T_AND_SELECTOR)
			return lexAttributeSelector
		} else if r == ':' {
			l.emit(T_AND_SELECTOR)
			return lexPseudoSelector
		} else if r == '#' {
			l.emit(T_AND_SELECTOR)
			return lexIdentifierSelector
		}
		return lexSelectors
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

// Dispath selector lexing method
func lexSelectors(l *Lexer) stateFn {
	var r = l.peek()

	// space between selector means descendant selector
	if tok := l.lastToken(); tok != nil && isSelector(tok.Type) {
		for r == ' ' {
			r = l.next()
		}
		l.backup()
		if r != '{' {
			l.emit(T_DESCENDANT_SELECTOR)
		}
	}

	// lex the first selector
	if unicode.IsLetter(r) {
		return lexTagNameSelector
	} else if r == '[' {
		return lexAttributeSelector
	} else if r == '.' {
		return lexClassSelector
	} else if r == ':' {
		return lexPseudoSelector
	} else if r == '&' {
		return lexParentSelector
	} else if r == '*' {
		return lexUniversalSelector
	} else if r == '#' {
		// for selector syntax like:
		//    '#{  }  {  }'
		//    '#{  }.something {  }'
		//    '#{  } .something {  }'
		//    '#{  }#myId {  }'
		if isInterpolationStartToken(r, l.peekMore(2)) {
			if tok := l.lastToken(); tok != nil && isSelector(tok.Type) {
				l.emit(T_CONCAT)
			}

			lexInterpolation(l, false)

			// find stop point of a selector.
			r = l.next()
			for !isSelectorStopToken(r) && !isDescendantSelectorSeparator(r) && isInterpolationStartToken(r, l.peekMore(2)) {
				r = l.next()
			}
			l.backup()

			// the suffix of the selector.
			var token = l.createToken(T_INTERPOLATION_SELECTOR)
			token.ContainsInterpolation = 1
			l.emitToken(token)
			// lext next selector
			return lexSelectors
		}
		return lexIdentifierSelector
	} else if r == '>' {
		return lexChildSelector
	} else if r == ',' {

		l.next()
		l.emit(T_COMMA)

		// lex next selector
		return lexSelectors

	} else if r == '+' {
		l.next()
		l.emit(T_PLUS)
		return lexSelectors
	} else if r == ' ' {
		for r == ' ' {
			r = l.next()
		}
		l.backup()
		l.ignore()
		return lexSelectors
	} else if r == '{' {
		return lexStatement
	} else {
		l.error("Unexpected token '%s' for selector.", r)
	}
	return nil
}

func lexTagNameSelector(l *Lexer) stateFn {
	var r = l.peek()
	var foundInterpolation = false
	if !unicode.IsLetter(r) && !isInterpolationStartToken(r, l.peekMore(2)) {
		l.error("Expecting letter token for tag name selector. got %s", r)
	}
	for {
		if isInterpolationStartToken(r, l.peekMore(2)) {
			lexInterpolation(l, true)
			foundInterpolation = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}
		r = l.next()
	}
	l.backup()

	if foundInterpolation {
		l.emit(T_INTERPOLATION_SELECTOR)
	} else {
		l.emit(T_TAGNAME_SELECTOR)
	}

	// predicate and inject the and selector for class name, identifier after the tagName
	switch r {
	case ':':
		l.emit(T_AND_SELECTOR)
		return lexPseudoSelector
	case '[':
		l.emit(T_AND_SELECTOR)
		return lexAttributeSelector
	case '#':
		if l.peekMore(2) != '{' {
			l.emit(T_AND_SELECTOR)
			return lexIdentifierSelector
		}
	case '.':
		l.emit(T_AND_SELECTOR)
		return lexClassSelector
	case '{':
		return lexStatement
	}
	return lexSelectors
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
