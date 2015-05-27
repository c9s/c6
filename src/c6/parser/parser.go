package parser

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "c6/ast"
import "c6/runtime"
import "path/filepath"
import "io/ioutil"
import "os"
import "fmt"

const (
	UnknownFileType uint = iota
	ScssFileType
	SassFileType
	EcssFileType
)

func debug(format string, args ...interface{}) {
	if debugParser {
		fmt.Printf(format+"\n", args...)
	}
}

func getFileTypeByExtension(extension string) uint {
	switch extension {
	case ".scss":
		return ScssFileType
	case ".sass":
		return SassFileType
	case ".ecss":
		return EcssFileType
	}
	return UnknownFileType
}

type Parser struct {
	GlobalContext *runtime.Context
	ContextStack  []runtime.Context

	File     string
	FileInfo os.FileInfo
	Content  string
	Input    chan *ast.Token

	// integer for counting token
	Pos         int
	RollbackPos int
	Tokens      []*ast.Token
}

func NewParser(context *runtime.Context) *Parser {
	return &Parser{
		GlobalContext: context,
		Input:         nil,
		Pos:           0,
		RollbackPos:   0,
	}
}

func (parser *Parser) ParseFile(path string) error {
	ext := filepath.Ext(path)
	filetype := getFileTypeByExtension(ext)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	parser.Content = string(data)
	parser.File = path

	switch filetype {
	case ScssFileType:
		parser.ParseScss(parser.Content)
		break
	}
	return nil
}

func (parser *Parser) backup() {
	parser.Pos--
}

func (parser *Parser) restore(pos int) {
	parser.Pos = pos
}

func (parser *Parser) remember() {
	parser.RollbackPos = parser.Pos
}

func (parser *Parser) rollback() {
	parser.Pos = parser.RollbackPos
}

func (parser *Parser) accept(tokenType ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok != nil && tok.Type == tokenType {
		return tok
	}
	parser.backup()
	return nil
}

func (parser *Parser) acceptAny(tokenTypes ...ast.TokenType) *ast.Token {
	var tok = parser.next()
	for _, tokType := range tokenTypes {
		if tok.Type == tokType {
			return tok
		}
	}
	parser.backup()
	return nil
}

func (parser *Parser) acceptAnyOf2(tokType1, tokType2 ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok.Type == tokType1 || tok.Type == tokType2 {
		return tok
	}
	parser.backup()
	return nil
}

func (parser *Parser) acceptAnyOf3(tokType1, tokType2, tokType3 ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok.Type == tokType1 || tok.Type == tokType2 || tok.Type == tokType3 {
		return tok
	}
	parser.backup()
	return nil
}

func (parser *Parser) expect(tokenType ast.TokenType) *ast.Token {
	var tok = parser.next()
	if tok != nil && tok.Type != tokenType {
		parser.backup()
		panic(SyntaxError{
			Reason:      tokenType.String(),
			ActualToken: tok,
			File:        parser.File,
		})
	}
	return tok
}

func (parser *Parser) next() *ast.Token {
	var p = parser.Pos
	parser.Pos++

	// if we've appended the token
	if p < len(parser.Tokens) {
		return parser.Tokens[p]
	}
	return nil
}

func (parser *Parser) peekBy(offset int) *ast.Token {
	var i = 0
	var tok *ast.Token = nil
	for offset > 0 {
		tok = parser.next()
		offset--
		i++
		if tok == nil {
			break
		}
	}
	parser.Pos -= i
	return tok
}

func (parser *Parser) advance() {
	parser.Pos++
}

func (parser *Parser) current() *ast.Token {
	return parser.Tokens[parser.Pos]
}

func (parser *Parser) peek() *ast.Token {
	if parser.Pos < len(parser.Tokens) {
		return parser.Tokens[parser.Pos]
	}
	return nil
}

func (parser *Parser) eof() bool {
	var tok = parser.next()
	parser.backup()
	return tok == nil
}
