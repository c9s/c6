package c6

import "io/ioutil"
import "path/filepath"
import "c6/ast"
import "strconv"

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
	Input chan *ast.Token

	// integer for counting token
	Pos         int
	RollbackPos int
	Tokens      []*ast.Token
}

func NewParser() *Parser {
	p := Parser{}
	p.Pos = 0
	p.Tokens = []*ast.Token{}
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

func (self *Parser) accept(tokenType ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type == tokenType {
		return tok
	}
	self.backup()
	return nil
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

	if token.Type == ast.T_IMPORT {
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
		case ast.T_TYPE_SELECTOR:
			sel := ast.TypeSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case ast.T_UNIVERSAL_SELECTOR:
			sel := ast.UniversalSelector{}
			ruleset.AppendSelector(sel)
		case ast.T_ID_SELECTOR:
			sel := ast.IdSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case ast.T_CLASS_SELECTOR:
			sel := ast.ClassSelector{tok.Str}
			ruleset.AppendSelector(sel)
		case ast.T_PARENT_SELECTOR:
			sel := ast.ParentSelector{parentRuleSet}
			ruleset.AppendSelector(sel)
		case ast.T_PSEUDO_SELECTOR:
			sel := ast.PseudoSelector{tok.Str, ""}
			if nextTok := parser.peek(); nextTok.Type == ast.T_LANG_CODE {
				sel.C = nextTok.Str
			}
			ruleset.AppendSelector(sel)
		case ast.T_ADJACENT_SELECTOR:
			ruleset.AppendSelector(ast.AdjacentSelector{})
		case ast.T_CHILD_SELECTOR:
			ruleset.AppendSelector(ast.ChildSelector{})
		case ast.T_DESCENDANT_SELECTOR:
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

func (parser *Parser) PushPropertyConstantStrings(property *ast.Property) {
	rightTok := parser.next()
	for rightTok.Type == ast.T_IDENT {
		property.AppendValue(ast.ConstantString{rightTok.Str, *rightTok})
		rightTok = parser.next()
	}
	parser.backup()
}

/**
This method returns objects with ast.Number interface

works for:

	'10'
	'10' 'px'
	'10' 'em'
	'0.2' 'em'
*/
func (parser *Parser) ReduceNumber() ast.Number {
	// the value token
	var value = parser.next()
	var tok2 = parser.peek()
	var number ast.Number
	if value.Type == ast.T_INTEGER {
		i, err := strconv.ParseInt(value.Str, 10, 64)
		if err != nil {
			panic(err)
		}
		number = ast.NewIntegerNumber(i)
	} else {

		f, err := strconv.ParseFloat(value.Str, 64)
		if err != nil {
			panic(err)
		}
		number = ast.NewFloatNumber(f)
	}

	if tok2.IsOneOfTypes([]ast.TokenType{ast.T_UNIT_PX, ast.T_UNIT_PT, ast.T_UNIT_CM, ast.T_UNIT_EM, ast.T_UNIT_MM, ast.T_UNIT_REM, ast.T_UNIT_DEG, ast.T_UNIT_PERCENT}) {
		number.SetUnit(int(tok2.Type))
	}
	return number
}

func (parser *Parser) ReduceFactor() {
	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_START {
		// skip the parent
		parser.accept(ast.T_PAREN_START)
		parser.ParseExpression()
		parser.accept(ast.T_PAREN_END)
		// _ = expr
	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {
		// reduce number

	}

}

func (parser *Parser) ReduceTerm() {
	parser.ReduceTerm()
}

/**

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};

*/
func (parser *Parser) ParseExpression() {
	parser.ReduceTerm()
	parser.acceptTypes([]ast.TokenType{ast.T_PLUS, ast.T_MINUS})
	parser.ReduceTerm()

}

/**
The returned Expression is an interface
*/
func (parser *Parser) ParsePropertyValue(parentRuleSet *ast.RuleSet, property *ast.Property) {
	tok := parser.next()

	// Check if the property value is an expression or
	if tok.Type == ast.T_IDENT {
		pos := parser.Pos

		// if it's a constant list
		rightTok := parser.next()
		for rightTok.Type == ast.T_IDENT {
			rightTok = parser.next()
		}
		// a constant list can be end with a semi colon or a end brace
		if rightTok.Type == ast.T_SEMICOLON || rightTok.Type == ast.T_BRACE_END {
			// OK, it's a constant list, restore to the original position (the position after the colon)
			parser.Pos = pos

			// push the values before semicolon into the property value list.
			parser.PushPropertyConstantStrings(property)
		} else {
			parser.Pos = pos
		}

		_ = pos
		_ = rightTok

		/*
			if right.Type == ast.T_SEMICOLON {
				parser.advance()
				// construct the normal property value since we meet semicolon
				return &ast.UnaryExpression{tok.Str, *tok}
			}
		*/
	} else {
		panic("unimplemented property value parsing...")
	}
}

func (parser *Parser) ParseDeclarationBlock(parentRuleSet *ast.RuleSet) *ast.DeclarationBlock {
	var declBlock = ast.DeclarationBlock{}

	var tok = parser.next() // should be '{'
	if tok.Type != ast.T_BRACE_START {
		panic(ParserError{"{", tok.Str})
	}

	tok = parser.next()
	for tok.Type != ast.T_BRACE_END {

		fmt.Printf("%+v\n", tok)

		if tok.Type == ast.T_PROPERTY_NAME {
			// skip T_COLON
			parser.next()

			var propertyName = ast.PropertyName{tok.Str, tok.ContainsInterpolation, *tok}
			var property = ast.Property{propertyName, []ast.Expression{}}

			tok := parser.peek()

			// Check if the property value is an expression or
			if tok.Type == ast.T_IDENT {

				parser.ParsePropertyValue(parentRuleSet, &property)
				property.AppendValue(ast.UnaryExpression{tok.Str, *tok})
				declBlock.Append(property)

				/*
					tok2 := parser.peek()
					if tok2.Type == ast.T_SEMICOLON {
						parser.advance()
						// construct the normal property value since we meet semicolon
						property.AppendValue(ast.UnaryExpression{tok.Str, *tok})
						declBlock.Append(property)

					} else if tok2.Type == ast.T_PAREN_START { // looks like a function

						panic("unimplemented function expression")
					}
				*/

			}

		} else if tok.IsSelector() {
			// parse subrule
			panic("subselector unimplemented")
		} else {
		}

		tok = parser.next()
	}

	return &declBlock
}

func (parser *Parser) ParseImportStatement() ast.Statement {
	// skip the ast.T_IMPORT token
	var tok = parser.next()

	// Create the import statement node
	var rule = ast.ImportStatement{}

	tok = parser.peek()
	// expecting url(..)
	if tok.Type == ast.T_IDENT {
		parser.advance()

		if tok.Str != "url" {
			panic("invalid function for @import rule.")
		}

		if tok = parser.next(); tok.Type != ast.T_PAREN_START {
			panic("expecting parenthesis after url")
		}

		tok = parser.next()
		rule.Url = ast.Url(tok.Str)

		if tok = parser.next(); tok.Type != ast.T_PAREN_END {
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
	if tok.Type == ast.T_MEDIA {
		parser.advance()
		rule.MediaList = append(rule.MediaList, tok.Str)
	}

	// must be ast.T_SEMICOLON
	tok = parser.next()
	if tok.Type != ast.T_SEMICOLON {
		panic(ParserError{";", tok.Str})
	}
	return &rule
}
