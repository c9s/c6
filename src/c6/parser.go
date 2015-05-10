package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "fmt"
import "c6/ast"
import "path/filepath"
import "io/ioutil"

const (
	UnknownFileType = iota
	ScssFileType
	SassFileType
	EcssFileType
)

type ParserError struct {
	ExpectingToken string
	ActualToken    string
}

const debugParser = false

func debug(format string, args ...interface{}) {
	if debugParser {
		fmt.Printf(format+"\n", args...)
	}
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Expecting '%s', but the actual token we got was '%s'.", e.ExpectingToken, e.ActualToken)
}

func getFileTypeByExtension(extension string) uint {
	switch extension {
	case "scss":
		return ScssFileType
	case "sass":
		return SassFileType
	case "ecss":
		return EcssFileType
	}
	return UnknownFileType
}

type Parser struct {
	Context *Context
	Input   chan *ast.Token

	// integer for counting token
	Pos         int
	RollbackPos int
	Tokens      []*ast.Token
}

func NewParser(context *Context) *Parser {
	return &Parser{context, nil, 0, 0, []*ast.Token{}}
}

func (parser *Parser) ParseFile(path string) error {
	ext := filepath.Ext(path)
	filetype := getFileTypeByExtension(ext)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var code string = string(data)
	switch filetype {
	case ScssFileType:
		parser.ParseScss(code)
		break
	}
	return nil
}

func (self *Parser) backup() {
	self.Pos--
}

func (self *Parser) restore(pos int) {
	self.Pos = pos
}

func (self *Parser) remember() {
	self.RollbackPos = self.Pos
}

func (self *Parser) rollback() {
	self.Pos = self.RollbackPos
}

func (self *Parser) accept(tokenType ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type == tokenType {
		return tok
	}
	self.backup()
	return nil
}

func (self *Parser) expect(tokenType ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type != tokenType {
		self.backup()
		panic(fmt.Errorf("Expecting %s, Got %s", tokenType, tok))
	}
	return tok
}

func (self *Parser) acceptTypes(types []ast.TokenType) bool {
	var p = self.Pos
	var match = true
	for _, tokType := range types {
		var tok = self.next()
		if tok.Type != tokType {
			match = false
			break
		}
	}
	// restore the position if it doesn't match
	if !match {
		self.Pos = p
	}
	return match
}

func (self *Parser) next() *ast.Token {
	var p = self.Pos
	self.Pos++
	if p < len(self.Tokens) {
		return self.Tokens[p]
	} else {
		if len(self.Tokens) > 1 {
			// get the last token
			var tok = self.Tokens[len(self.Tokens)-1]
			if tok == nil {
				return nil
			}
		}
		token := <-self.Input
		self.Tokens = append(self.Tokens, token)
		return token
	}
	return nil
}

func (self *Parser) peekBy(offset int) *ast.Token {
	var i = 0
	var tok *ast.Token = nil
	for offset > 0 && tok != nil {
		tok = self.next()
		offset--
		i++
	}
	self.Pos -= i
	return tok
}

func (self *Parser) advance() {
	self.Pos++
}

func (self *Parser) current() *ast.Token {
	return self.Tokens[self.Pos]
}

func (self *Parser) peek() *ast.Token {
	if self.Pos < len(self.Tokens) {
		return self.Tokens[self.Pos]
	}
	token := <-self.Input
	self.Tokens = append(self.Tokens, token)
	return token
}

func (self *Parser) isSelector() bool {
	var tok = self.peek()
	if tok.Type == ast.T_ID_SELECTOR ||
		tok.Type == ast.T_TYPE_SELECTOR ||
		tok.Type == ast.T_CLASS_SELECTOR ||
		tok.Type == ast.T_PSEUDO_SELECTOR ||
		tok.Type == ast.T_PARENT_SELECTOR {
		return true
	} else if tok.Type == ast.T_BRACKET_LEFT {
		return true
	}
	return false
}

func (self *Parser) eof() bool {
	var tok = self.next()
	self.backup()
	return tok == nil
}

func (parser *Parser) ParseScss(code string) []ast.Statement {
	l := NewLexerWithString(code)
	l.run()
	parser.Input = l.getOutput()
	return parser.ParseStatements()
}
