package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "io/ioutil"
import "unicode/utf8"
import "strings"
import "fmt"
import "unicode"
import "c6/ast"

type tokenChannel chan *ast.Token

const TOKEN_CHANNEL_BUFFER = 1024

const EOF = -1

const DEBUG_EMIT = true

type Lexer struct {
	// lex input
	Input string

	// current buffer offset
	Offset int

	// the offset where token starts
	Start int

	// byte width of the current rune (utf8 character has more than one bytes)
	// The width will be updated by 'next()` method
	// `backup()` use Width to go back to the last offset.
	Width int

	// After the next() is called, the original width is backed up in
	// LastWidth
	LastWidth int

	// rollback offset for token
	RollbackOffset int

	// current lexer file
	File string

	// current lexer state
	State stateFn

	// current line number of the input
	Line int

	// the token output channel
	Output chan *ast.Token

	Tokens []ast.Token
}

func (l *Lexer) lastToken() *ast.Token {
	if len(l.Tokens) > 0 {
		return &l.Tokens[len(l.Tokens)-1]
	}
	return nil
}

/**
Create a lexer object with bytes
*/
func NewLexerWithBytes(data []byte) *Lexer {
	l := &Lexer{
		File:   "{anonymous}",
		Offset: 0,
		Line:   0,
		Input:  string(data),
	}
	return l
}

/**
Create a lexer object with string
*/
func NewLexerWithString(body string) *Lexer {
	return &Lexer{
		File:   "{anonymous}",
		Offset: 0,
		Line:   0,
		Input:  body,
	}
}

/**
Create a lexer object with file path

TODO: detect encoding here
*/
func NewLexerWithFile(file string) (*Lexer, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &Lexer{
		File:   file,
		Offset: 0,
		Line:   0,
		Input:  string(data),
	}, nil
}

func (l *Lexer) getOutput() chan *ast.Token {
	if l.Output != nil {
		return l.Output
	}
	l.Output = make(chan *ast.Token, TOKEN_CHANNEL_BUFFER)
	return l.Output
}

// remember the current offset, can be rolled back by using the `rollback`
// method
func (l *Lexer) remember() int {
	l.RollbackOffset = l.Offset
	return l.Offset
}

// rollback reset the offset to the backup point (this is a rune-wise
// rollback)
func (l *Lexer) rollback() {
	l.Offset = l.RollbackOffset
}

func (l *Lexer) acceptAndEmit(valid string, tokenType ast.TokenType) bool {
	if l.accept(valid) {
		l.emit(tokenType)
		return true
	}
	return false
}

func (l *Lexer) expect(valid string) {
	if !l.accept(valid) {
		panic(fmt.Errorf("Expecting %s at %d", valid, l.Offset))
	}
}

// test the next character, if it's not matched, go back to the original
// offset.
// Note, this method only match the first character
func (l *Lexer) accept(valid string) bool {
	var r rune = l.next()
	if strings.IndexRune(valid, r) >= 0 {
		return true
	}
	l.backup()
	return false
}

// Accept letter runes continuously
// Return true if there are some letters.
// Retunr false if there is no letter.
func (l *Lexer) acceptLetters() bool {
	var r rune = l.next()
	for unicode.IsLetter(r) {
		r = l.next()
	}
	l.backup()
	return l.Offset > l.Start
}

// Accept letter|digits runes continuously
// Return true if there are some letters.
// Retunr false if there is no letter.
func (l *Lexer) acceptLettersAndDigits() bool {
	var r rune = l.next()
	for unicode.IsLetter(r) || unicode.IsDigit(r) {
		r = l.next()
	}
	l.backup()
	return l.Offset > l.Start
}

// Return the current token string but not consume it
func (l *Lexer) current() string {
	if l.Offset >= len(l.Input) {
		return ""
	}
	return l.Input[l.Start:l.Offset]
}

// Return the length of the current token
func (l *Lexer) length() int {
	return l.Offset - l.Start
}

// If there are remaining tokens
func (l *Lexer) remaining() bool {
	return l.Offset+1 < len(l.Input)
}

// next returns the next rune in the input.
func (l *Lexer) next() (r rune) {
	if l.Offset >= len(l.Input) {
		l.LastWidth = l.Width
		l.Width = 0
		return EOF
	}
	l.LastWidth = l.Width
	r, l.Width = utf8.DecodeRuneInString(l.Input[l.Offset:])
	l.Offset += l.Width
	return r
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backup() {
	l.Offset -= l.Width
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backupByWidth(w int) {
	l.Offset -= w
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() (r rune) {
	r = l.next()
	l.backup()
	return r
}

// advance offset by specific width
func (l *Lexer) advance(w int) {
	l.Offset += w
}

// peek more characters
// peekBy(1) == peek()
func (l *Lexer) peekBy(p int) (r rune) {
	var w = 0
	for i := p; i > 0; i-- {
		r = l.next()
		w += l.Width
	}
	l.Offset -= w
	return r
}

func (l *Lexer) take() string {
	return l.Input[l.Start:l.Offset]
}

func (l *Lexer) emitToken(token *ast.Token) {
	if DEBUG_EMIT {
		fmt.Printf("emit: %+v\n", token)
	}

	l.Tokens = append(l.Tokens, *token)
	l.Output <- token
	l.Start = l.Offset
}

func (l *Lexer) createTokenWith0Offset(tokenType ast.TokenType) *ast.Token {
	var token = ast.Token{
		Type: tokenType,
		Str:  "",
		Pos:  l.Start,
		Line: l.Line,
	}
	return &token
}

func (l *Lexer) createToken(tokenType ast.TokenType) *ast.Token {
	/*
		if l.Offset > len(l.Input) {
			panic(fmt.Sprintf("out of range at '%s': start:%d, offset:%d, length: %d", l.Input[l.Start:], l.Start, l.Offset, len(l.Input)))
		}
	*/
	var token = ast.Token{
		Type: tokenType,
		Str:  l.Input[l.Start:l.Offset],
		Pos:  l.Start,
		Line: l.Line,
	}
	return &token
}

/*
Emit a token to the channel

emit(ast.T_SEMICOLON)

emit(ast.T_PSEUDO_SELECTOR, true) // contains interpolation
*/
func (l *Lexer) emit(tokenType ast.TokenType, params ...bool) {
	token := l.createToken(tokenType)
	if len(params) > 0 && params[0] {
		token.ContainsInterpolation = true
	}
	l.emitToken(token)
}

// lookahead a string til {string}
func (l *Lexer) lookaheadTil(stop string) string {
	l.remember()
	for {
		var r = l.next()
		if strings.Contains(stop, string(r)) || r == EOF {
			break
		}
	}
	var str = l.take()
	l.rollback()
	return str
}

func (l *Lexer) til(str string) {
	for {
		var r = l.next()
		if strings.Contains(str, string(r)) || r == EOF {
			break
		}
	}
	l.backup()
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.Start = l.Offset
}

// lookahead match method to match a string
func (l *Lexer) match(str string) bool {
	var r rune
	var width = 0
	for _, sc := range str {
		r = l.next()
		width += l.Width
		if sc != r {
			// rollback
			l.Offset -= width
			return false
		}
	}
	return true
}

func (l *Lexer) matchKeywordMap(keywords ast.KeywordTokenMap) ast.TokenType {
	for str, tokType := range keywords {
		l.remember()
		if l.match(str) {
			var r = l.peek()
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
				// try next one
				l.rollback()
				continue
			}
			l.emit(tokType)
			return tokType
		}
	}
	return 0
}

func (l *Lexer) precedeStartOffset() bool {
	return l.Offset > l.Start
}

/*
ignore space characters

return true if there is space
*/
func (l *Lexer) ignoreSpaces() int {
	var space = 0
	for {
		var r rune = l.peek()
		if r == '\n' {
			space++
			l.Line++
			l.next()
		} else if r == ' ' || r == '\t' || r == '\r' {
			space++
			l.next()
		} else {
			break
		}
	}
	// Update the token start offset to latest offset
	l.Start = l.Offset
	return space
}

func (l *Lexer) dispatchFn(fn stateFn) stateFn {
	for l.State = fn; l.State != nil; {
		fn := l.State(l)
		if fn != nil {
			l.State = fn
		} else {
			break
		}
	}
	return l.State
}

func (l *Lexer) dump() {
	fmt.Printf("Lexer: %+v\n", l)
}

func (l *Lexer) runFrom(fn stateFn) {
	if l.Output == nil {
		l.Output = make(tokenChannel, TOKEN_CHANNEL_BUFFER)
	}
	l.dispatchFn(fn)
	l.Output <- nil
}

func (l *Lexer) run() {
	if l.Output == nil {
		l.Output = make(tokenChannel, TOKEN_CHANNEL_BUFFER)
	}
	l.dispatchFn(lexStart)
	l.Output <- nil
}

func (l *Lexer) close() {
	if l.Output != nil {
		close(l.Output)
	}
}
