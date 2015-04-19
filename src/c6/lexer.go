package c6

import "io/ioutil"
import "unicode/utf8"
import "strings"
import "fmt"

type tokenChannel chan *Token

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
	Output chan *Token

	Tokens []Token
}

func (l *Lexer) lastToken() Token {
	return l.Tokens[len(self.Tokens)-1]
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

func (l *Lexer) getOutput() chan *Token {
	if l.Output != nil {
		return l.Output
	}
	l.Output = make(chan *Token, TOKEN_CHANNEL_BUFFER)
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
func (l *Lexer) peekMore(p int) (r rune) {
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

// emit a token to the channel
func (l *Lexer) emit(tokenType TokenType) {
	token := Token{
		Type: tokenType,
		Str:  l.Input[l.Start:l.Offset],
		Pos:  l.Start,
		Line: l.Line,
	}
	if DEBUG_EMIT {
		fmt.Println("emit", token)
	}

	l.Tokens = append(l.Tokens, token)
	l.Output <- &token
	l.Start = l.Offset
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
			l.Offset -= width
			return false
		}
	}
	return true
}

func (l *Lexer) ignoreSpaces() {
	for {
		var r rune = l.peek()
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			l.next()
		} else {
			break
		}
	}
	// Update the token start offset to latest offset
	l.Start = l.Offset
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

/*
func (self *Lexer) peek() {
	var p = self.Offset
	if self.State == StateRoot {
		if self.Input[p] == '.' {
			p++
			for {
				var c = self.Input[p]
				if c == ' ' || c == '{' {
					break
				}
				if !unicode.IsLetter(c) && c != '-' {
					break
				}
			}
		}
	}
}
*/

/*
func (self *Lexer) lexSelector() *Token {
	return self.lexClassSelector()
}
*/

/*
func (self *Lexer) lexClassSelector() *Token {
	var p = self.Offset
	if self.Input[p] == '.' {
		p++

		// TODO: Prevent p to overflow here
		for {
			var c = self.Input[p]
			// if it's the end of a .class selector
			if c == ' ' || c == '{' {
				return &Token{
					Type: TokenClassSelector,
					Str:  self.Input[self.Offset : p-1],
					Pos:  self.Offset,
				}
				break
			}
			if !unicode.IsLetter(c) && c != '-' {
				// Raise error here
				break
			}
			p++
		}
	}
	return nil
}
*/
