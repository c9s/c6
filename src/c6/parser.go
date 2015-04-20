package c6

import "io/ioutil"
import "path/filepath"

var fileAstMap map[string]interface{} = map[string]interface{}{}

const (
	UnknownFileType = iota
	ScssFileType
	SassFileType
)

func getFileTypeByExtension(extension string) uint {
	switch extension {
	case "scss":
		return ScssFileType
	case "sass":
		return SassFileType
	}
	return UnknownFileType
}

type Parser struct {
	Input chan *Token

	// integer for counting token
	Pos    int
	Tokens []Token
}

func NewParser() *Parser {
	p := Parser{}
	p.Pos = 0
	p.Tokens = []Token{}
	return &p
}

func (parser *Parser) parseFile(path string) error {
	ext := filepath.Ext(path)
	filetype := getFileTypeByExtension(ext)
	_ = filetype
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var code string = string(data)
	_ = code
	return nil
}

func (self *Parser) next() *Token {
	var p = self.Pos
	self.Pos++
	if p < len(self.Tokens) {
		return &self.Tokens[p]
	} else if token := <-self.Input; token != nil {
		self.Tokens = append(self.Tokens, *token)
		return token
	}
	return nil
}

func (self *Parser) peek() *Token {
	if self.Pos < len(self.Tokens) {
		return &self.Tokens[self.Pos]
	}

	if token := <-self.Input; token != nil {
		self.Tokens = append(self.Tokens, *token)
		return token
	}
	return nil
}

func (self *Parser) parseScss(code string) {
	l := NewLexerWithString(code)
	l.run()
	self.Input = l.getOutput()
}

func (self *Parser) parseSass(code string) {

}
