package c6

import "io/ioutil"
import "path/filepath"
import "c6/ast"

// import "errors"
import "fmt"

var fileAstMap map[string]interface{} = map[string]interface{}{}

const (
	UnknownFileType = iota
	ScssFileType
	SassFileType
)

type ParserError struct {
	ExpectingToken string
	ActualToken    string
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
	}
	return UnknownFileType
}

type Parser struct {
	Input chan *Token

	// integer for counting token
	Pos         int
	RollbackPos int
	Tokens      []*Token
}

func NewParser() *Parser {
	p := Parser{}
	p.Pos = 0
	p.Tokens = []*Token{}
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

func (self *Parser) matchByTypes(types []TokenType) bool {
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

func (self *Parser) next() *Token {
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

func (self *Parser) peekBy(offset int) *Token {
	if self.Pos+offset < len(self.Tokens) {
		return self.Tokens[self.Pos+offset]
	}
	token := <-self.Input
	for token != nil {
		self.Tokens = append(self.Tokens, token)
		if self.Pos+offset < len(self.Tokens) {
			return self.Tokens[self.Pos+offset]
		}
		token = <-self.Input
	}
	return nil
}

func (self *Parser) advance() {
	self.Pos++
}

func (self *Parser) current() *Token {
	return self.Tokens[self.Pos]
}

func (self *Parser) peek() *Token {
	if self.Pos < len(self.Tokens) {
		return self.Tokens[self.Pos]
	}
	token := <-self.Input
	self.Tokens = append(self.Tokens, token)
	return token
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

func (self *Parser) eof() bool {
	var tok = self.next()
	self.backup()
	return tok == nil
}

func (parser *Parser) parseScss(code string) *ast.Block {
	l := NewLexerWithString(code)
	l.run()
	parser.Input = l.getOutput()

	block := ast.Block{}
	for !parser.eof() {
		stm := parser.ParseStatement(nil)
		if stm != nil {
			block.AppendStatement(stm)
		}
	}
	return &block
}

/*
Statement := RuleSet | At-Rule | Mixin-Statement | FunctionStatement


At-Rule := '@' T_IDENT ';'

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
func (parser *Parser) ParseStatement(parentRuleSet *ast.RuleSet) ast.Statement {
	var token = parser.peek()

	if token.Type == T_IMPORT {
		return parser.ParseImportStatement()
	} else if token.IsSelector() {
		return parser.ParseRuleSet(parentRuleSet)
	}
	return nil
}

func (parser *Parser) ParseRuleSet(parentRuleSet *ast.RuleSet) ast.Statement {
	var ruleset = ast.RuleSet{}
	var tok = parser.next()

	for tok.IsSelector() {

		switch tok.Type {
		case T_TYPE_SELECTOR:
			sel := ast.TypeSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case T_UNIVERSAL_SELECTOR:
			sel := ast.UniversalSelector{}
			ruleset.AppendSelector(sel)
		case T_ID_SELECTOR:
			sel := ast.IdSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case T_CLASS_SELECTOR:
			sel := ast.ClassSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case T_PARENT_SELECTOR:
			sel := ast.ParentSelector{parentRuleSet}
			ruleset.AppendSelector(sel)
		case T_PSEUDO_SELECTOR:
			sel := ast.PseudoSelector{tok.Str, ""}
			if nextTok := parser.peek(); nextTok.Type == T_LANG_CODE {
				sel.C = nextTok.Str
			}
			ruleset.AppendSelector(sel)
		case T_ADJACENT_SELECTOR:
			ruleset.AppendSelector(ast.AdjacentSelector{})
		case T_CHILD_SELECTOR:
			ruleset.AppendSelector(ast.ChildSelector{})
		case T_DESCENDANT_SELECTOR:
			ruleset.AppendSelector(ast.DescendantSelector{})
		default:
			panic(fmt.Errorf("Unexpected selector token: %+v", tok))
		}
		tok = parser.next()
	}
	parser.backup()

	// parse declaration block
	ruleset.DeclarationBlock = parser.ParseDeclarationBlock(&ruleset)
	return &ruleset
}

func (parser *Parser) ParseDeclarationBlock(parentRuleSet *ast.RuleSet) *ast.DeclarationBlock {
	var tok = parser.next() // should be '{'
	if tok.Type != T_BRACE_START {
		panic(ParserError{"(", tok.Str})
	}

	tok = parser.next()
	for tok.Type != T_BRACE_END {
		tok = parser.next()
	}

	return &ast.DeclarationBlock{}
}

func (parser *Parser) ParseImportStatement() ast.Statement {
	// skip the T_IMPORT token
	var tok = parser.next()

	// Create the import statement node
	var rule = ast.ImportStatement{}

	tok = parser.peek()
	// expecting url(..)
	if tok.Type == T_IDENT {
		parser.advance()

		if tok.Str != "url" {
			panic("invalid function for @import rule.")
		}

		if tok = parser.next(); tok.Type != T_PAREN_START {
			panic("expecting parenthesis after url")
		}

		tok = parser.next()
		rule.Url = ast.Url(tok.Str)

		if tok = parser.next(); tok.Type != T_PAREN_END {
			panic("expecting parenthesis after url")
		}

	} else if tok.IsString() {
		parser.advance()
		rule.Url = ast.RelativeUrl(tok.Str)
	}

	/*
		TODO: parse media query for something like:

		@import url(color.css) screen and (color);
		@import url('landscape.css') screen and (orientation:landscape);
		@import url("bluish.css") projection, tv;
	*/
	tok = parser.peek()
	if tok.Type == T_MEDIA {
		parser.advance()
		rule.MediaList = append(rule.MediaList, tok.Str)
	}

	// must be T_SEMICOLON
	tok = parser.next()
	if tok.Type != T_SEMICOLON {
		panic(ParserError{";", tok.Str})
	}
	return &rule
}
