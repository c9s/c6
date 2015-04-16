package c6

import "unicode"

func IsNewLine(c rune) bool {
	return c == '\n'
}

func IsClassSelectorChar(r rune) bool {
	return r == '.' || unicode.IsLetter(r) || r == '-' || r == '_'
}

func IsIDSelectorChar(r rune) bool {
	return r == '#' || unicode.IsLetter(r) || r == '-' || r == '_'
}

func IsPropertyNameChar(c rune) bool {
	return unicode.IsLetter(c) || c == '-'
}

func IsPropertySepChar(c rune) bool {
	return c == ':'
}

func IsComma(c rune) bool {
	return c == ';'
}
