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
	Pos         int
	RollbackPos int
	Tokens      []Token
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

func (self *Parser) backup() {
	self.Pos--
}

func (self *Parser) remember() {
	self.RollbackPos = self.Pos
}

func (self *Parser) rollback() {
	self.Pos = self.RollbackPos
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

func (self *Parser) peekBy(offset int) *Token {
	if self.Pos+offset < len(self.Tokens) {
		return &self.Tokens[self.Pos+offset]
	}
	token := <-self.Input
	for token != nil {
		self.Tokens = append(self.Tokens, *token)
		if self.Pos+offset < len(self.Tokens) {
			return &self.Tokens[self.Pos+offset]
		}
		token = <-self.Input
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

func (self *Parser) isSelector() bool {
	var tok = self.peek()
	if tok.Type == T_ID_SELECTOR ||
		tok.Type == T_TYPE_SELECTOR ||
		tok.Type == T_CLASS_SELECTOR ||
		tok.Type == T_PSEUDO_SELECTOR ||
		tok.Type == T_PARENT_SELECTOR {
		return true
	} else if tok.Type == T_BRACKET_LEFT {
		return true
	}
	return false
}

func (self *Parser) parseScss(code string) {
	l := NewLexerWithString(code)
	l.run()
	self.Input = l.getOutput()
	/*
		if self.isSelector() {
			rule := Rule{}
			_ = rule
		}
	*/
}

/*
Statement := RuleSet | At-Rule | Mixin-Statement | FunctionStatement


At-Rule := '@' T_IDENTITY ';'

RuleSet := Rule | RuleSet

Property := PropertyName ':' PropertyValue

PropertyValue := Expr ';'
		       | ConstantList ';'

ConstantList := T_CONSTANT | T_CONSTANT ConstantList

SelectorList := Selector | Selector ',' SelectorList

Rule := SelectorList '{'
			RuleSet
		'}'

Expr := Expr '+' Expr
      | Expr '-' Expr
      | Expr '*' Expr
      | Expr '/' Expr
	  | '(' Expr ')'
	  | Scalar
	  | Variable

Variable := T_VARIABLE

Scalar := T_NUMBER | T_NUMBER Unit

Unit := T_UNIT_PX | T_UNIT_PT | T_UNIT_EM | T_UNIT_PERCENT | T_UNIT_DEG
*/

func (self *Parser) parseStatement() {

}

func (self *Parser) parseAtRule() {

}

func (self *Parser) parseRule() {

}
