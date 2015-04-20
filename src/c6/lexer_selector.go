package c6

import "unicode"

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
			return lexClassSelector
		} else if r == '[' {
			l.emit(T_AND_SELECTOR)
			return lexAttributeSelector
		}
		return lexSelector
	}
	return nil
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
		return lexSelector
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

func lexChildSelector(l *Lexer) stateFn {
	var r = l.next()
	if r == '>' {
		l.emit(T_GT)
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
			l.ignore()
			lexLang(l)
			r = l.next()
			if r != ')' {
				l.error("Unexpected token '%s' for pseudo lang selector", r)
			}
			l.ignore()
		}
		return lexSelector
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
		return lexSelector
	}
	l.error("Unexpected token '%s' for universal selector.", r)
	return nil
}

// Dispath selector lexing method
func lexSelector(l *Lexer) stateFn {
	var r = l.peek()

	if unicode.IsLetter(r) {
		return lexTagNameSelector
	} else if r == '[' {
		return lexAttributeSelector
	} else if r == '.' {
		return lexClassSelector
	} else if r == '#' && l.peekMore(2) == '{' {
		lexInterpolation(l)
		return lexSelector
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
		l.emit(T_PLUS)
		return lexSelector
	} else if r == ' ' {
		for r == ' ' {
			r = l.next()
		}
		l.backup()
		l.ignore()
		return lexSelector
	} else if r == '{' {
		return lexStatement
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
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
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
