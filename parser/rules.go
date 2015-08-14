package parser

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/c9s/c6/ast"
	"github.com/c9s/c6/lexer"
	"github.com/c9s/c6/runtime"
)

var HttpUrlPattern = regexp.MustCompile("^https?://")
var AbsoluteUrlPattern = regexp.MustCompile("^[a-zA-Z]+?://")

func (parser *Parser) ReadFile(file string) error {
	f, err := ast.NewFile(file)
	if err != nil {
		return err
	}
	data, err := f.ReadFile()
	if err != nil {
		return err
	}
	parser.File = f
	parser.Content = string(data)
	return nil
}

func (parser *Parser) ParseScssFile(file string) (*ast.StmtList, error) {
	err := parser.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// XXX: this seems to copy the whole string, we should avoid this.
	l := lexer.NewLexerWithString(parser.Content)
	parser.Input = l.TokenStream()

	// Run lexer concurrently
	go l.Run()

	// consume the tokens from the input channel of the lexer
	// TODO: use concurrent method to consume the inputs, we also need to
	// benchmark this when the file is large. Don't need to consider small files because small files
	// can always be compiled fast (less than 500 millisecond).
	var tok *ast.Token = nil
	for tok = <-parser.Input; tok != nil; tok = <-parser.Input {
		parser.Tokens = append(parser.Tokens, tok)
	}
	l.Close()
	return parser.ParseStmts(), nil
}

func (parser *Parser) ParseScss(code string) *ast.StmtList {
	l := lexer.NewLexerWithString(code)
	parser.Input = l.TokenStream()

	// Run lexer concurrently
	go l.Run()

	var tok *ast.Token = nil
	for tok = <-parser.Input; tok != nil; tok = <-parser.Input {
		parser.Tokens = append(parser.Tokens, tok)
	}
	l.Close()
	return parser.ParseStmts()
}

/*
ParseBlock method allows root level statements, which does not allow css properties.
*/
func (parser *Parser) ParseBlock() *ast.Block {
	debug("ParseBlock")
	parser.expect(ast.T_BRACE_OPEN)
	var block = ast.NewBlock()
	block.Stmts = parser.ParseStmts()
	parser.expect(ast.T_BRACE_CLOSE)
	return block
}

func (parser *Parser) ParseStmts() *ast.StmtList {
	var stmts = new(ast.StmtList)
	// stop at t_brace end
	for !parser.eof() {
		if stm := parser.ParseStmt(); stm != nil {
			*stmts = append(*stmts, stm)
		} else {
			break
		}
	}
	return stmts
}

func (parser *Parser) ParseStmt() ast.Stmt {
	var token = parser.peek()

	if token == nil {
		return nil
	}

	switch token.Type {
	case ast.T_IMPORT:
		return parser.ParseImportStmt()
	case ast.T_CHARSET:
		return parser.ParseCharsetStmt()
	case ast.T_MEDIA:
		return parser.ParseMediaQueryStmt()
	case ast.T_MIXIN:
		return parser.ParseMixinStmt()
	case ast.T_FUNCTION:
		return parser.ParseFunctionDeclaration()
	case ast.T_FONT_FACE:
		return parser.ParseFontFaceStmt()
	case ast.T_INCLUDE:
		return parser.ParseIncludeStmt()
	case ast.T_VARIABLE:
		return parser.ParseAssignStmt()
	case ast.T_RETURN:
		return parser.ParseReturnStmt()
	case ast.T_IF:
		return parser.ParseIfStmt()
	case ast.T_EXTEND:
		return parser.ParseExtendStmt()
	case ast.T_FOR:
		return parser.ParseForStmt()
	case ast.T_WHILE:
		return parser.ParseWhileStmt()
	case ast.T_CONTENT:
		return parser.ParseContentStmt()
	case ast.T_ERROR, ast.T_WARN, ast.T_INFO:
		return parser.ParseLogStmt()
	case ast.T_BRACKET_CLOSE:
		return nil
	}

	if token.IsSelector() {

		return parser.ParseRuleSet()

	} else {

		// panic(fmt.Errorf("statement parse failed, unknown token", parser.peek()))

	}
	return nil
}

func (parser *Parser) ParseIfStmt() ast.Stmt {
	parser.expect(ast.T_IF)
	condition := parser.ParseCondition()
	if condition == nil {
		panic("if statement syntax error")
	}

	var block = parser.ParseDeclBlock()
	var stm = ast.NewIfStmt(condition, block)

	// TODO: OptimizeIfStmt(...)

	// If these is more else if statement
	var tok = parser.peek()
	for tok != nil && tok.Type == ast.T_ELSE_IF {
		parser.advance()

		// XXX: handle error here
		var condition = parser.ParseCondition()
		var elseifblock = parser.ParseDeclBlock()
		var elseIfStm = ast.NewIfStmt(condition, elseifblock)
		stm.AppendElseIf(elseIfStm)
		tok = parser.peek()
	}

	tok = parser.peek()
	if tok != nil && tok.Type == ast.T_ELSE {
		parser.advance()

		// XXX: handle error here
		if elseBlock := parser.ParseDeclBlock(); elseBlock != nil {
			stm.ElseBlock = elseBlock
		} else {
			panic(SyntaxError{
				Reason:      "Expecting declaration block { ... }",
				ActualToken: parser.peek(),
				File:        parser.File,
			})
		}
	}

	return stm
}

/*
The operator precedence is described here

@see http://introcs.cs.princeton.edu/java/11precedence/
*/
func (parser *Parser) ParseCondition() ast.Expr {
	debug("ParseCondition")

	// Boolean 'Not'
	if tok := parser.accept(ast.T_LOGICAL_NOT); tok != nil {
		var logicexpr = parser.ParseLogicExpr()
		return ast.NewUnaryExpr(ast.NewOpWithToken(tok), logicexpr)
	}
	return parser.ParseLogicExpr()
}

func (parser *Parser) ParseLogicExpr() ast.Expr {
	debug("ParseLogicExpr")
	var expr = parser.ParseLogicANDExpr()
	for tok := parser.accept(ast.T_LOGICAL_OR); tok != nil; tok = parser.accept(ast.T_LOGICAL_OR) {
		if subexpr := parser.ParseLogicANDExpr(); subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
	}
	return expr
}

func (parser *Parser) ParseLogicANDExpr() ast.Expr {
	debug("ParseLogicANDExpr")

	var expr = parser.ParseComparisonExpr()
	for tok := parser.accept(ast.T_LOGICAL_AND); tok != nil; tok = parser.accept(ast.T_LOGICAL_AND) {
		if subexpr := parser.ParseComparisonExpr(); subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
	}
	return expr
}

func (parser *Parser) ParseComparisonExpr() ast.Expr {
	debug("ParseComparisonExpr")

	var expr ast.Expr = nil
	if parser.accept(ast.T_PAREN_OPEN) != nil {
		expr = parser.ParseLogicExpr()
		parser.expect(ast.T_PAREN_CLOSE)
	} else {
		expr = parser.ParseExpr(false)
	}

	var tok = parser.peek()
	for tok != nil && tok.IsComparisonOperator() {
		parser.advance()
		if subexpr := parser.ParseExpr(false); subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
		tok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseSimpleSelector(parentRuleSet *ast.RuleSet) ast.Selector {
	debug("ParseSimpleSelector")

	var tok = parser.next()
	if tok == nil {
		return nil
	}

	switch tok.Type {

	case ast.T_TYPE_SELECTOR:

		return ast.NewTypeSelectorWithToken(tok)

	case ast.T_UNIVERSAL_SELECTOR:

		return ast.NewUniversalSelectorWithToken(tok)

	case ast.T_ID_SELECTOR:

		return ast.NewIdSelectorWithToken(tok)

	case ast.T_CLASS_SELECTOR:

		return ast.NewClassSelectorWithToken(tok)

	case ast.T_PARENT_SELECTOR:

		return ast.NewParentSelectorWithToken(parentRuleSet, tok)

	case ast.T_FUNCTIONAL_PSEUDO:

		var sel = ast.NewFunctionalPseudoSelectorWithToken(tok)
		parser.expect(ast.T_PAREN_OPEN)

		var tok2 = parser.next()
		for tok2 != nil && tok2.Type != ast.T_PAREN_CLOSE {
			// TODO: parse pseudo expression
			tok2 = parser.next()
		}
		parser.backup()
		parser.expect(ast.T_PAREN_CLOSE)

		return sel

	case ast.T_PSEUDO_SELECTOR:

		return ast.NewPseudoSelectorWithToken(tok)

	// Attribute selector parsing
	case ast.T_BRACKET_OPEN:
		var attrName = parser.expect(ast.T_ATTRIBUTE_NAME)
		var tok2 = parser.next()
		if tok2.IsAttributeMatchOperator() {

			var tok3 = parser.next()
			var sel = ast.NewAttributeSelector(attrName, tok2, tok3)
			parser.expect(ast.T_BRACKET_CLOSE)
			return sel

		} else if tok2.Type == ast.T_BRACKET_CLOSE {

			return ast.NewAttributeSelectorNameOnly(attrName)

		} else {

			panic(SyntaxError{
				Reason:      "Unexpected token",
				ActualToken: tok2,
				File:        parser.File,
			})
		}
	}
	parser.backup()
	return nil
}

func (parser *Parser) ParseCompoundSelector(parentRuleSet *ast.RuleSet) *ast.CompoundSelector {
	var sels = ast.NewCompoundSelector()
	for {
		if sel := parser.ParseSimpleSelector(parentRuleSet); sel != nil {
			sels.Append(sel)
		} else {
			break
		}
	}
	if sels.Length() > 0 {
		return sels
	}
	return nil
}

func (parser *Parser) ParseComplexSelector(parentRuleSet *ast.RuleSet) *ast.ComplexSelector {
	debug("ParseComplexSelector")

	var sel = parser.ParseCompoundSelector(parentRuleSet)
	if sel == nil {
		return nil
	}

	var complexSel = ast.NewComplexSelector(sel)

	for {
		var tok = parser.next()
		if tok == nil {
			return complexSel
		}

		var comb ast.Combinator

		// peek the combinator token
		switch tok.Type {

		case ast.T_ADJACENT_SIBLING_COMBINATOR:

			comb = ast.NewAdjacentCombinatorWithToken(tok)

		case ast.T_CHILD_COMBINATOR:

			comb = ast.NewChildCombinatorWithToken(tok)

		case ast.T_DESCENDANT_COMBINATOR:

			comb = ast.NewDescendantCombinatorWithToken(tok)

		case ast.T_GENERAL_SIBLING_COMBINATOR:

			comb = ast.NewGeneralSiblingCombinatorWithToken(tok)

		default:
			parser.backup()
			return complexSel
		}

		if sel := parser.ParseCompoundSelector(parentRuleSet); sel != nil {

			complexSel.AppendCompoundSelector(comb, sel)

		} else {

			panic(SyntaxError{
				Reason:      "Expecting a selector after the combinator.",
				ActualToken: parser.peek(),
				File:        parser.File,
			})

		}
	}
}

func (parser *Parser) ParseSelectorList() *ast.ComplexSelectorList {
	debug("ParseSelectorList")

	var parentRuleSet = parser.GlobalContext.TopRuleSet()

	var complexSelectorList = &ast.ComplexSelectorList{}

	if complexSelector := parser.ParseComplexSelector(parentRuleSet); complexSelector != nil {
		complexSelectorList.Append(complexSelector)
	} else {
		return nil
	}

	// if there is more comma
	for parser.accept(ast.T_COMMA) != nil {
		if complexSelector := parser.ParseComplexSelector(parentRuleSet); complexSelector != nil {
			complexSelectorList.Append(complexSelector)
		} else {
			break
		}
	}
	return complexSelectorList
}

func (parser *Parser) ParseExtendStmt() ast.Stmt {
	parser.expect(ast.T_EXTEND)
	var stm = ast.NewExtendStmt()
	stm.Selectors = parser.ParseSelectorList()
	parser.expect(ast.T_SEMICOLON)
	return stm
}

func (parser *Parser) ParseRuleSet() ast.Stmt {
	var ruleset = ast.NewRuleSet()
	ruleset.Selectors = parser.ParseSelectorList()

	parser.GlobalContext.PushRuleSet(ruleset)
	ruleset.Block = parser.ParseDeclBlock()
	parser.GlobalContext.PopRuleSet()

	return ruleset
}

func (parser *Parser) ParseBoolean() ast.Expr {
	if tok := parser.acceptAnyOf2(ast.T_TRUE, ast.T_FALSE); tok != nil {
		return ast.NewBooleanWithToken(tok)
	}
	return nil
}

func (parser *Parser) ParseNumber() ast.Expr {
	var pos = parser.Pos
	debug("ParseNumber at %d", parser.Pos)

	// the number token
	var tok = parser.next()
	debug("ParseNumber => next: %s", tok)

	var negative = false

	if tok.Type == ast.T_MINUS {
		tok = parser.next()
		negative = true
	} else if tok.Type == ast.T_PLUS {
		tok = parser.next()
		negative = false
	}

	var val float64
	var tok2 = parser.peek()

	if tok.Type == ast.T_INTEGER {

		i, err := strconv.ParseInt(tok.Str, 10, 64)
		if err != nil {
			panic(err)
		}
		if negative {
			i = -i
		}
		val = float64(i)

	} else if tok.Type == ast.T_FLOAT {

		f, err := strconv.ParseFloat(tok.Str, 64)
		if err != nil {
			panic(err)
		}
		if negative {
			f = -f
		}
		val = f

	} else {
		// unknown token
		parser.restore(pos)
		return nil
	}

	if tok2.IsUnit() {
		// consume the unit token
		parser.next()
		return ast.NewNumber(val, ast.NewUnitWithToken(tok2), tok)
	}
	return ast.NewNumber(val, nil, tok)
}

/*
This parse method looks up the argument in the function declaration and convert
the keyword arguments into the argument list. this way, we can handle the function call
in a simple way - push/pop the stack.

@param fcall ast.FunctionCall The current parsing function call ast node
@return arguments []Expr
*/
func (parser *Parser) ParseKeywordArguments(fcall *ast.FunctionCall) *ast.FunctionCallArguments {
	// look up function declaration
	var item, ok = parser.GlobalContext.Functions.Get(fcall.Ident.Str)
	if !ok {
		panic("Undefined function " + fcall.Ident.Str)
	}
	var fun = item.(*ast.Function)

	var arguments = ast.FunctionCallArguments{}
	for tok := parser.accept(ast.T_VARIABLE); tok != nil; tok = parser.accept(ast.T_VARIABLE) {
		var argdef = fun.ArgumentList.Lookup(tok.Str)
		if argdef == nil {
			panic("Undefined function argument: " + tok.Str)
		}

		parser.expect(ast.T_COLON)
		var argExpr = parser.ParseExpr(false)
		parser.accept(ast.T_COMMA)

		var arg = ast.NewFunctionCallArgument(argExpr)
		arg.ArgumentDefineReference = argdef

		arguments = append(arguments, arg)
	}

	// sort arguments by ArgumentDefineReference
	arguments.Sort()
	return &arguments
}

func (parser *Parser) ParseFunctionCall() *ast.FunctionCall {
	var identTok = parser.next()

	debug("ParseFunctionCall => next: %s", identTok)

	var fcall = ast.NewFunctionCallWithToken(identTok)

	parser.expect(ast.T_PAREN_OPEN)

	var tok = parser.peek()
	var tok2 = parser.peekBy(2)
	if tok.Type == ast.T_VARIABLE && tok2.Type == ast.T_COLON {

		fcall.Arguments = *parser.ParseKeywordArguments(fcall)

	} else {
		for tok.Type != ast.T_PAREN_CLOSE {

			if arg := parser.ParseFactor(); arg != nil {
				fcall.AppendArgument(arg)
				debug("ParseFunctionCall => arg: %+v", arg)
			} else {
				break
			}

			if parser.accept(ast.T_COMMA) != nil {
				tok = parser.peek()
				continue
			}
		}
	}
	parser.expect(ast.T_PAREN_CLOSE)
	return fcall
}

func (parser *Parser) ParseIdent() *ast.Ident {
	var tok = parser.next()
	if tok.Type != ast.T_IDENT {
		panic("Invalid token for ident.")
	}
	return ast.NewIdentWithToken(tok)
}

/**
The ParseFactor must return an Expr interface compatible object
*/
func (parser *Parser) ParseFactor() ast.Expr {
	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_OPEN {

		parser.expect(ast.T_PAREN_OPEN)
		var expr = parser.ParseExpr(true)
		parser.expect(ast.T_PAREN_CLOSE)
		return expr

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	} else if tok.Type == ast.T_QQ_STRING {

		parser.advance()
		return ast.NewStringWithQuote('"', tok)

	} else if tok.Type == ast.T_Q_STRING {

		parser.advance()
		return ast.NewStringWithQuote('\'', tok)

	} else if tok.Type == ast.T_UNQUOTE_STRING {

		parser.advance()
		return ast.NewStringWithQuote(0, tok)

	} else if tok.Type == ast.T_TRUE {

		parser.advance()
		return ast.NewBooleanTrue(tok)

	} else if tok.Type == ast.T_FALSE {

		parser.advance()
		return ast.NewBooleanFalse(tok)

	} else if tok.Type == ast.T_NULL {

		parser.advance()
		return ast.NewNullWithToken(tok)

	} else if tok.Type == ast.T_FUNCTION_NAME {

		var fcall = parser.ParseFunctionCall()
		return ast.Expr(fcall)

	} else if tok.Type == ast.T_VARIABLE {

		return parser.ParseVariable()

	} else if tok.Type == ast.T_IDENT {

		var tok2 = parser.peekBy(2)
		if tok2 != nil && tok2.Type == ast.T_PAREN_OPEN {
			return parser.ParseFunctionCall()
		}

		parser.advance()
		return ast.NewStringWithToken(tok)

	} else if tok.Type == ast.T_HEX_COLOR {

		parser.advance()
		return ast.NewHexColorFromToken(tok)

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		return parser.ParseNumber()

	}
	return nil
}

func (parser *Parser) ParseTerm() ast.Expr {
	var pos = parser.Pos
	var factor = parser.ParseFactor()
	if factor == nil {
		parser.restore(pos)
		return nil
	}

	// see if the next token is '*' or '/'
	if tok := parser.acceptAnyOf2(ast.T_MUL, ast.T_DIV); tok != nil {
		if term := parser.ParseTerm(); term != nil {
			return ast.NewBinaryExpr(ast.NewOpWithToken(tok), factor, term, false)
		} else {
			panic(SyntaxError{
				Reason:      "Expecting term after '*' or '/'",
				ActualToken: parser.peek(),
				File:        parser.File,
			})
		}
	}
	return factor
}

/**

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};

*/
func (parser *Parser) ParseExpr(inParenthesis bool) ast.Expr {
	var pos = parser.Pos

	// plus or minus. This creates an unary expression that holds the later term.
	// this is for:  +3 or -4
	var expr ast.Expr = nil

	if tok := parser.acceptAnyOf2(ast.T_PLUS, ast.T_MINUS); tok != nil {
		if term := parser.ParseTerm(); term != nil {
			expr = ast.NewUnaryExpr(ast.NewOpWithToken(tok), term)

			if uexpr, ok := expr.(*ast.UnaryExpr); ok {

				// if it's evaluatable just return the evaluated value.
				if val, ok := runtime.ReduceExpr(uexpr, parser.GlobalContext); ok {
					expr = ast.Expr(val)
				}
			}
		} else {
			parser.restore(pos)
			return nil
		}
	} else {
		expr = parser.ParseTerm()
	}

	if expr == nil {
		debug("ParseExpr failed, got %+v, restoring to %d", expr, pos)
		parser.restore(pos)
		return nil
	}

	var rightTok = parser.peek()
	for rightTok.Type == ast.T_PLUS || rightTok.Type == ast.T_MINUS || rightTok.Type == ast.T_LITERAL_CONCAT {
		// accept plus or minus
		parser.advance()

		if rightTerm := parser.ParseTerm(); rightTerm != nil {
			// XXX: check parenthesis
			var bexpr = ast.NewBinaryExpr(ast.NewOpWithToken(rightTok), expr, rightTerm, inParenthesis)

			if val, ok := runtime.ReduceExpr(bexpr, parser.GlobalContext); ok {

				expr = ast.Expr(val)

			} else {
				// wrap the existing expression with the new binary expression object
				expr = ast.Expr(bexpr)
			}
		} else {
			panic(SyntaxError{
				Reason:      "Expecting term on the right side",
				ActualToken: parser.peek(),
				File:        parser.File,
			})
		}
		rightTok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseMap() ast.Expr {
	var pos = parser.Pos
	var tok = parser.accept(ast.T_PAREN_OPEN)
	// since it's not started with '(', it's not map
	if tok == nil {
		parser.restore(pos)
		return nil
	}

	var mapval = ast.NewMap()

	// TODO: check and report Map syntax error
	tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var keyExpr = parser.ParseExpr(false)
		if keyExpr == nil {
			parser.restore(pos)
			return nil
		}

		if parser.accept(ast.T_COLON) == nil {
			parser.restore(pos)
			return nil
		}

		var valueExpr = parser.ParseExpr(false)
		if valueExpr == nil {
			parser.restore(pos)
			return nil
		}

		// register the map value
		mapval.Set(keyExpr, valueExpr)
		parser.accept(ast.T_COMMA)
		tok = parser.peek()
	}
	if parser.accept(ast.T_PAREN_CLOSE) == nil {
		return nil
	}
	return mapval
}

func (parser *Parser) ParseString() ast.Expr {
	if tok := parser.accept(ast.T_QQ_STRING); tok != nil {

		return ast.NewStringWithQuote('"', tok)

	} else if tok := parser.accept(ast.T_Q_STRING); tok != nil {

		return ast.NewStringWithQuote('\'', tok)

	} else if tok := parser.accept(ast.T_UNQUOTE_STRING); tok != nil {

		return ast.NewStringWithQuote(0, tok)

	} else if tok := parser.accept(ast.T_IDENT); tok != nil {

		return ast.NewStringWithToken(tok)

	}

	var tok = parser.peek()
	if tok.Type == ast.T_INTERPOLATION_START {
		return parser.ParseInterp()
	}
	return nil
}

func (parser *Parser) ParseInterp() ast.Expr {
	var startTok = parser.expect(ast.T_INTERPOLATION_START)
	var innerExpr = parser.ParseExpr(true)
	var endTok = parser.expect(ast.T_INTERPOLATION_END)
	return ast.NewInterpolation(innerExpr, startTok, endTok)
}

func (parser *Parser) ParseValueStrict() ast.Expr {
	var pos = parser.Pos

	if tok := parser.accept(ast.T_PAREN_OPEN); tok != nil {
		if mapValue := parser.ParseMap(); mapValue != nil {
			return mapValue
		}
		parser.restore(pos)

		if listValue := parser.ParseList(); listValue != nil {
			return listValue
		}
		parser.restore(pos)
	}
	// Reduce the expression
	return parser.ParseExpr(false)
}

/*
Parse string literal expression (literal concat with interpolation)
*/
func (parser *Parser) ParseLiteralExpr() ast.Expr {
	if expr := parser.ParseExpr(false); expr != nil {
		for tok := parser.accept(ast.T_LITERAL_CONCAT); tok != nil; tok = parser.accept(ast.T_LITERAL_CONCAT) {
			var rightExpr = parser.ParseExpr(false)
			if rightExpr == nil {
				panic(SyntaxError{
					Reason:      "Expecting expression or ident after the literal concat operator.",
					ActualToken: parser.peek(),
					File:        parser.File,
				})
			}
			expr = ast.NewLiteralConcat(expr, rightExpr)
		}

		// Check if the expression is reduce-able
		// For now, division looks like CSS slash at the first level, should be string.
		if runtime.CanReduceExpr(expr) {
			if reducedExpr, ok := runtime.ReduceExpr(expr, parser.GlobalContext); ok {
				return reducedExpr
			}
		} else {
			// Return expression as css slash syntax string
			// TODO: re-visit here later
			return runtime.EvaluateExpr(expr, parser.GlobalContext)
		}

		// if we can't evaluate the value, just return the expression tree
		return expr
	}
	return nil
}

/**
Parse Value loosely

This parse method allows space/comma separated tokens of list.

To parse mixin argument or function argument, we only allow comma-separated list inside the parenthesis.

@param stopTokType ast.TokenType

The stop token is used from variable assignment expression,
We expect ';' semicolon at the end of expression to avoid the ambiguity of list, map and expression.
*/
func (parser *Parser) ParseValue(stopTokType ast.TokenType) ast.Expr {
	var pos = parser.Pos

	// try parse map
	debug("ParseMap")
	if mapValue := parser.ParseMap(); mapValue != nil {
		var tok = parser.peek()
		if stopTokType == 0 || tok.Type == stopTokType {
			debug("OK Map Meet Stop Token")
			return mapValue
		}
	}
	debug("Map parse failed, restoring to %d", pos)
	parser.restore(pos)

	debug("Trying List")
	if listValue := parser.ParseList(); listValue != nil {
		var tok = parser.peek()
		if stopTokType == 0 || tok.Type == stopTokType {
			debug("OK List: %+v", listValue)
			return listValue
		}
	}

	debug("List parse failed, restoring to %d", pos)
	parser.restore(pos)

	debug("ParseLiteralExpr trying", pos)
	return parser.ParseLiteralExpr()
}

func (parser *Parser) ParseList() ast.Expr {
	debug("ParseList at %d", parser.Pos)
	var pos = parser.Pos
	if list := parser.ParseCommaSepList(); list != nil {
		return list
	}
	parser.restore(pos)
	return nil
}

func (parser *Parser) ParseCommaSepList() ast.Expr {
	debug("ParseCommaSepList at %d", parser.Pos)
	var list = ast.NewCommaSepList()

	var tok = parser.peek()
	for tok.Type != ast.T_COMMA && tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {

		// when the syntax start with a '(', it could be a list or map.
		if tok.Type == ast.T_PAREN_OPEN {

			parser.next()
			if sublist := parser.ParseCommaSepList(); sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			}
			// allow empty list here
			parser.expect(ast.T_PAREN_CLOSE)

		} else {
			var sublist = parser.ParseSpaceSepList()
			if sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			} else {
				break
			}
		}

		if parser.accept(ast.T_COMMA) == nil {
			break
		}
		tok = parser.peek()
	}

	debug("Returning comma-separated list: (%+v)", list)

	if list.Len() == 0 {

		return nil

	} else if list.Len() == 1 {

		return list.Exprs[0]

	}
	return list
}

func (parser *Parser) ParseVariable() *ast.Variable {
	if tok := parser.accept(ast.T_VARIABLE); tok != nil {
		return ast.NewVariableWithToken(tok)
	}
	return nil
}

func (parser *Parser) ParseAssignStmt() ast.Stmt {
	var variable = parser.ParseVariable()

	// skip ":", T_COLON token
	parser.expect(ast.T_COLON)

	// Expecting semicolon at the end of the statement
	var valExpr = parser.ParseValue(ast.T_SEMICOLON)
	if valExpr == nil {
		panic(SyntaxError{
			Reason:      "Expecting value after variable assignment.",
			ActualToken: parser.peek(),
			File:        parser.File,
		})
	}

	// Optimize the expression only when it's an expression
	// TODO: for expression inside a map or list we should also optmise them too
	if bexpr, ok := valExpr.(ast.BinaryExpr); ok {
		if reducedExpr, ok := runtime.ReduceExpr(bexpr, parser.GlobalContext); ok {
			valExpr = reducedExpr
		}
	} else if uexpr, ok := valExpr.(ast.UnaryExpr); ok {
		if reducedExpr, ok := runtime.ReduceExpr(uexpr, parser.GlobalContext); ok {
			valExpr = reducedExpr
		}
	}

	// FIXME
	// Even we can visit the variable assignment in the AST visitors but if we
	// could save the information, we can reduce the effort for the visitors.
	/*
		if currentBlock := parser.GlobalContext.CurrentBlock(); currentBlock != nil {
			currentBlock.GetSymTable().Set(variable.Name, valExpr)
		} else {
			panic("nil block")
		}
	*/

	var stm = ast.NewAssignStmt(variable, valExpr)
	parser.ParseFlags(stm)
	parser.accept(ast.T_SEMICOLON)
	return stm
}

/*
ParseFlags requires a variable assignment.
*/
func (parser *Parser) ParseFlags(stm *ast.AssignStmt) {
	var tok = parser.peek()
	for tok.IsFlagKeyword() {
		parser.next()

		switch tok.Type {
		case ast.T_FLAG_DEFAULT:
			stm.Default = true
		case ast.T_FLAG_OPTIONAL:
			stm.Optional = true
		case ast.T_FLAG_IMPORTANT:
			stm.Important = true
		case ast.T_FLAG_GLOBAL:
			stm.Global = true
		}
		tok = parser.peek()
	}
}

func (parser *Parser) ParseSpaceSepList() ast.Expr {
	debug("ParseSpaceSepList at %d", parser.Pos)

	var list = ast.NewSpaceSepList()
	list.Separator = " "

	if tok := parser.accept(ast.T_PAREN_OPEN); tok != nil {
		if sublist := parser.ParseCommaSepList(); sublist != nil {
			list.Append(sublist)
		}
		parser.expect(ast.T_PAREN_CLOSE)
	}

	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {
		var subexpr = parser.ParseExpr(true)
		if subexpr != nil {
			debug("Parsed Expr: %+v", subexpr)
			list.Append(subexpr)
		} else {
			break
		}
		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			break
		}
	}
	debug("Returning space-sep list: %+v", list)
	if list.Len() == 0 {
		return nil
	} else if list.Len() == 1 {
		return list.Exprs[0]
	} else if list.Len() > 1 {
		return list
	}
	return nil
}

/**
We treat the property value section as a list value, which is separated by ',' or ' '
*/
func (parser *Parser) ParsePropertyValue(parentRuleSet *ast.RuleSet, property *ast.Property) *ast.List {
	debug("ParsePropertyValue")
	// var tok = parser.peek()
	var list = ast.NewSpaceSepList()

	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {
		var sublist = parser.ParseList()
		if sublist != nil {
			list.Append(sublist)
			debug("ParsePropertyValue list: %+v", list)
		} else {
			break
		}
		tok = parser.peek()
	}

	return list
}

func (parser *Parser) ParsePropertyName() ast.Expr {
	var ident = parser.ParsePropertyNameToken()
	if ident == nil {
		return nil
	}

	var tok = parser.peek()
	for tok.Type == ast.T_LITERAL_CONCAT {
		parser.next()
		_ = parser.ParsePropertyNameToken()
		tok = parser.peek()
	}
	parser.expect(ast.T_COLON)
	return ident // TODO: new literal concat ast
}

func (parser *Parser) ParsePropertyNameToken() ast.Expr {
	if tok := parser.accept(ast.T_PROPERTY_NAME_TOKEN); tok != nil {
		return ast.NewIdentWithToken(tok)
	}

	var tok = parser.peek()
	if tok.Type == ast.T_INTERPOLATION_START {
		return parser.ParseInterpolation()
	}
	return nil
}

func (parser *Parser) ParseInterpolation() ast.Expr {
	debug("ParseInterpolation")
	var startToken *ast.Token
	if startToken = parser.accept(ast.T_INTERPOLATION_START); startToken == nil {
		return nil
	}
	var expr = parser.ParseExpr(true)
	var endToken = parser.expect(ast.T_INTERPOLATION_END)
	return ast.NewInterpolation(expr, startToken, endToken)
}

func (parser *Parser) ParseDeclaration() ast.Stmt {
	return nil
}

func (parser *Parser) ParseDeclBlock() *ast.DeclBlock {
	var parentRuleSet = parser.GlobalContext.TopRuleSet()
	var declBlock = ast.NewDeclBlock(parentRuleSet)

	parser.expect(ast.T_BRACE_OPEN)

	var tok = parser.peek()
	for tok != nil && tok.Type != ast.T_BRACE_CLOSE {
		if propertyName := parser.ParsePropertyName(); propertyName != nil {
			var property = ast.NewProperty(tok)

			if valueList := parser.ParsePropertyValue(parentRuleSet, property); valueList != nil {
				// property.Values = valueList
			}
			declBlock.Append(property)

			var tok2 = parser.peek()

			// if nested property found
			if tok2.Type == ast.T_BRACE_OPEN {
				// TODO: merge them back to current block
				parser.ParseDeclBlock()
			}

			if parser.accept(ast.T_SEMICOLON) == nil {
				if tok3 := parser.peek(); tok3.Type == ast.T_BRACE_CLOSE {
					// normal break
					break
				} else {
					panic("missing semicolon after the property value.")
				}
			}

		} else if stm := parser.ParseStmt(); stm != nil {

		} else {
			panic(fmt.Errorf("Parse failed at token %s", tok))
		}
		tok = parser.peek()

	}
	parser.expect(ast.T_BRACE_CLOSE)
	return declBlock
}

func (parser *Parser) ParseCharsetStmt() ast.Stmt {
	parser.expect(ast.T_CHARSET)
	var tok = parser.next()
	var stm = ast.NewCharsetStmtWithToken(tok)
	parser.expect(ast.T_SEMICOLON)
	return stm
}

/*
	Media Query Syntax:
	https://developer.mozilla.org/en-US/docs/Web/Guide/CSS/Media_queries
*/
func (parser *Parser) ParseMediaQueryStmt() ast.Stmt {
	// expect the '@media' token
	var stm = ast.NewMediaQueryStmt()
	parser.expect(ast.T_MEDIA)
	if list := parser.ParseMediaQueryList(); list != nil {
		stm.MediaQueryList = list
	}
	stm.Block = parser.ParseDeclBlock()
	return stm
}

func (parser *Parser) ParseMediaQueryList() *ast.MediaQueryList {
	var query = parser.ParseMediaQuery()
	if query == nil {
		return nil
	}

	var queries = &ast.MediaQueryList{query}
	for parser.accept(ast.T_COMMA) != nil {
		if query := parser.ParseMediaQuery(); query != nil {
			queries.Append(query)
		}
	}
	return queries
}

/*
This method parses media type first, then expecting more that on media
expressions.

media_query: [[only | not]? <media_type> [ and <expression> ]*]
  | <expression> [ and <expression> ]*
expression: ( <media_feature> [: <value>]? )

Specification: http://dev.w3.org/csswg/mediaqueries-3
*/
func (parser *Parser) ParseMediaQuery() *ast.MediaQuery {

	// the leading media type is optional
	var mediaType = parser.ParseMediaType()
	if mediaType != nil {
		// Check if there is an expression after the media type.
		var tok = parser.peek()
		if tok.Type != ast.T_LOGICAL_AND {
			return ast.NewMediaQuery(mediaType, nil)
		}
		parser.advance() // skip the and operator token
	}

	// parse the media expression after the media type.
	var mediaExpr = parser.ParseMediaQueryExpr()
	if mediaExpr == nil {
		if mediaType == nil {
			return nil
		}
		return ast.NewMediaQuery(mediaType, mediaExpr)
	}

	// @media query only allows AND operator here..
	for tok := parser.accept(ast.T_LOGICAL_AND); tok != nil; tok = parser.accept(ast.T_LOGICAL_AND) {
		// parse another mediq query expression
		var expr2 = parser.ParseMediaQueryExpr()
		mediaExpr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), mediaExpr, expr2, false)
	}
	return ast.NewMediaQuery(mediaType, mediaExpr)
}

/*
ParseMediaType returns Ident Node or UnaryExpr as ast.Expr
*/
func (parser *Parser) ParseMediaType() *ast.MediaType {
	if tok := parser.acceptAnyOf2(ast.T_LOGICAL_NOT, ast.T_ONLY); tok != nil {
		var mediaType = parser.expect(ast.T_IDENT)
		return ast.NewMediaType(ast.NewUnaryExpr(ast.NewOpWithToken(tok), mediaType))
	}

	var tok = parser.peek()
	if tok.Type == ast.T_PAREN_OPEN {
		// the begining of the media expression
		return nil
	}

	var expr = parser.ParseExpr(false)
	if expr != nil {
		return ast.NewMediaType(expr)
	}

	// parse media type fail
	return nil
}

/*
An media query expression must start with a '(' and ends with ')'
*/
func (parser *Parser) ParseMediaQueryExpr() ast.Expr {

	// it's not an media query expression
	if openTok := parser.accept(ast.T_PAREN_OPEN); openTok != nil {
		var featureExpr = parser.ParseExpr(false)
		var feature = ast.NewMediaFeature(featureExpr, nil)

		// if the next token is a colon, then we expect a feature value
		// after the colon.
		if tok := parser.accept(ast.T_COLON); tok != nil {
			feature.Value = parser.ParseExpr(false)
		}
		var closeTok = parser.expect(ast.T_PAREN_CLOSE)
		feature.Open = openTok
		feature.Close = closeTok
		return feature
	}
	return nil
}

func (parser *Parser) ParseWhileStmt() ast.Stmt {
	parser.expect(ast.T_WHILE)
	var condition = parser.ParseCondition()
	var block = parser.ParseDeclBlock()
	return ast.NewWhileStmt(condition, block)
}

/*
Parse the SASS @for statement.

	@for $var from <start> to <end> {  }

	@for $var from <start> through <end> {  }

@see http://sass-lang.com/documentation/file.SASS_REFERENCE.html#_10
*/
func (parser *Parser) ParseForStmt() ast.Stmt {
	parser.expect(ast.T_FOR)

	// get the variable token
	var variable = parser.ParseVariable()
	var stm = ast.NewForStmt(variable)

	if parser.accept(ast.T_FOR_FROM) != nil {

		var fromExpr = parser.ParseExpr(true)
		if reducedExpr, ok := runtime.ReduceExpr(fromExpr, parser.GlobalContext); ok {
			fromExpr = reducedExpr
		}
		stm.From = fromExpr

		// "through" or "to"
		var tok = parser.next()

		if tok.Type != ast.T_FOR_THROUGH && tok.Type != ast.T_FOR_TO {
			panic(SyntaxError{
				Reason:      "Expecting 'through' or 'to' of range syntax.",
				ActualToken: tok,
				File:        parser.File,
			})
		}

		var endExpr = parser.ParseExpr(true)
		if reducedExpr, ok := runtime.ReduceExpr(endExpr, parser.GlobalContext); ok {
			endExpr = reducedExpr
		}

		if tok.Type == ast.T_FOR_THROUGH {

			stm.Through = endExpr

		} else if tok.Type == ast.T_FOR_TO {

			stm.To = endExpr

		}

	} else if parser.accept(ast.T_FOR_IN) != nil {

		var fromExpr = parser.ParseExpr(true)
		if reducedExpr, ok := runtime.ReduceExpr(fromExpr, parser.GlobalContext); ok {
			fromExpr = reducedExpr
		}
		stm.From = fromExpr

		parser.expect(ast.T_RANGE)

		var endExpr = parser.ParseExpr(true)
		if reducedExpr, ok := runtime.ReduceExpr(endExpr, parser.GlobalContext); ok {
			endExpr = reducedExpr
		}

		stm.To = endExpr
	}

	if b := parser.ParseDeclBlock(); b != nil {
		stm.Block = b
	} else {
		panic("The @for statement expecting block after the range syntax")
	}
	return stm
}

/*
The @import syntax is described here:

@see CSS2.1 http://www.w3.org/TR/CSS2/cascade.html#at-import

@see https://developer.mozilla.org/en-US/docs/Web/CSS/@import
*/
func (parser *Parser) ParseImportStmt() ast.Stmt {
	// skip the ast.T_IMPORT token
	parser.expect(ast.T_IMPORT)

	// Create the import statement node
	var stm = ast.NewImportStmt()

	var tok = parser.peek()

	// if it's url(..)
	if tok.Type == ast.T_FUNCTION_NAME {
		parser.advance()
		parser.expect(ast.T_PAREN_OPEN)

		var urlTok = parser.acceptAnyOf3(ast.T_QQ_STRING, ast.T_Q_STRING, ast.T_UNQUOTE_STRING)
		if urlTok == nil {
			panic(SyntaxError{
				Reason:      "Expecting url string in the url() function expression",
				ActualToken: parser.peek(),
				File:        parser.File,
			})
		}

		if HttpUrlPattern.MatchString(urlTok.Str) {
			stm.Url = ast.AbsoluteUrl(urlTok.Str)
		} else {
			stm.Url = ast.RelativeUrl(urlTok.Str)
		}

		parser.expect(ast.T_PAREN_CLOSE)

	} else if tok.IsString() {

		parser.advance()

		// Relative url for CSS
		if strings.HasSuffix(tok.Str, ".css") {

			stm.Url = ast.StringUrl(tok.Str)

		} else if AbsoluteUrlPattern.MatchString(tok.Str) {

			stm.Url = ast.AbsoluteUrl(tok.Str)

		} else if strings.HasSuffix(tok.Str, ".scss") {

			stm.Url = ast.ScssImportUrl(tok.Str)

		} else {
			// check scss import url by file system
			if parser.File != nil {
				panic("Unknown scss file to detect import path.")
			}

			var importPath = tok.Str
			var fi, err = os.Stat(importPath)
			if err != nil {
				panic(err)
			}

			// go find the _index.scss if it's a local directory
			if fi.Mode().IsDir() {

				importPath = importPath + string(filepath.Separator) + "_index.scss"

			} else {
				var dirname = filepath.Dir(importPath)
				var basename = filepath.Base(importPath)
				if dirname == "." {
					importPath = "_" + basename + ".scss"
				} else {
					importPath = dirname + string(filepath.Separator) + "_" + basename + ".scss"
				}
			}

			if _, err := os.Stat(importPath); err != nil {
				panic(err)
			}
			stm.Url = ast.ScssImportUrl(importPath)

			if _, ok := parser.GlobalContext.ImportedPath[importPath]; !ok {
				// Set imported path to true
				parser.GlobalContext.ImportedPath[importPath] = true

				// parse the imported file using the same context
				var subparser = NewParser(parser.GlobalContext)
				var stmts, err = subparser.ParseScssFile(importPath)
				if err != nil {
					panic(err)
				}

				// Cache the compiled statement in the AST, so we can include the compiled AST nodes when watch mode is enabled
				// if it's not changed.
				ASTCache[importPath] = stmts

				// For root @import and nested ruleset @import:
				//
				// 1. For nested rules, we need to merge the rulesets and assignment to the parent ruleset
				//    we can expand the the statements in the parsing-time.
				//
				// 2. for root level statements, we need to merge the statements into the global block.
				//    for symbal table, we also need to merge them
				//
				// note that the parse method might push the statements to global block, we should avoid that.
				var currentBlock = parser.GlobalContext.CurrentBlock()
				currentBlock.MergeStmts(stmts)
			}

		}

	} else {
		panic(fmt.Errorf("Unexpected token: %s", tok))
	}

	if mediaQueryList := parser.ParseMediaQueryList(); mediaQueryList != nil {
		stm.MediaQueryList = *mediaQueryList
	}

	// must be ast.T_SEMICOLON at the end
	parser.expect(ast.T_SEMICOLON)
	return stm
}

func (parser *Parser) ParseReturnStmt() ast.Stmt {
	var returnTok = parser.expect(ast.T_RETURN)
	var valueExpr = parser.ParseExpr(true)
	parser.expect(ast.T_SEMICOLON)
	return ast.NewReturnStmtWithToken(returnTok, valueExpr)
}

func (parser *Parser) ParseFunctionDeclaration() ast.Stmt {
	parser.expect(ast.T_FUNCTION)
	var identTok = parser.expect(ast.T_FUNCTION_NAME)
	var args = parser.ParseFunctionPrototype()
	var fun = ast.NewFunctionWithToken(identTok)
	fun.ArgumentList = args
	fun.Block = parser.ParseBlock()
	parser.GlobalContext.Functions.Set(identTok.Str, fun)
	return fun
}

func (parser *Parser) ParseMixinStmt() ast.Stmt {
	var mixinTok = parser.expect(ast.T_MIXIN)
	var stm = ast.NewMixinStmtWithToken(mixinTok)

	var tok = parser.next()

	// Mixin without parameters
	if tok.Type == ast.T_IDENT {

		stm.Ident = tok

	} else if tok.Type == ast.T_FUNCTION_NAME {

		stm.Ident = tok
		stm.ArgumentList = parser.ParseFunctionPrototype()

	} else {

		panic("Syntax error")

	}

	stm.Block = parser.ParseDeclBlock()
	parser.GlobalContext.Mixins.Set(stm.Ident.Str, stm)
	return stm
}

func (parser *Parser) ParseFunctionPrototypeArgument() *ast.Argument {
	debug("ParseFunctionPrototypeArgument")

	var varTok *ast.Token = nil
	if varTok = parser.accept(ast.T_VARIABLE); varTok == nil {
		return nil
	}

	if arg := ast.NewArgumentWithToken(varTok); arg != nil {
		if parser.accept(ast.T_COLON) != nil {
			arg.DefaultValue = parser.ParseValueStrict()
		}
		return arg
	}
	return nil
}

func (parser *Parser) ParseFunctionPrototype() *ast.ArgumentList {
	debug("ParseFunctionPrototype")

	var args = ast.NewArgumentList()

	parser.expect(ast.T_PAREN_OPEN)
	var tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var arg *ast.Argument = nil
		if arg = parser.ParseFunctionPrototypeArgument(); arg != nil {
			args.Add(arg)
		} else {
			// if fail
			break
		}
		if tok = parser.accept(ast.T_COMMA); tok != nil {
			continue
		} else if tok = parser.accept(ast.T_VARIABLE_LENGTH_ARGUMENTS); tok != nil {
			arg.VariableLength = true
			break
		} else {
			break
		}
	}
	parser.expect(ast.T_PAREN_CLOSE)
	return args
}

func (parser *Parser) ParseFunctionCallArgument() *ast.Argument {
	debug("ParseFunctionCallArgument")

	var varTok *ast.Token = nil
	if varTok = parser.accept(ast.T_VARIABLE); varTok == nil {
		return nil
	}

	var arg = ast.NewArgumentWithToken(varTok)

	if parser.accept(ast.T_COLON) != nil {
		arg.DefaultValue = parser.ParseValueStrict()
	}
	return arg
}

func (parser *Parser) ParseFunctionCallArguments() *ast.ArgumentList {
	debug("ParseFunctionCallArguments")

	var args = ast.NewArgumentList()

	parser.expect(ast.T_PAREN_OPEN)
	var tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var arg *ast.Argument = nil
		if arg = parser.ParseFunctionCallArgument(); arg != nil {
			args.Add(arg)
		} else {
			// if fail
			break
		}
		if tok = parser.accept(ast.T_COMMA); tok != nil {
			continue
		} else if tok = parser.accept(ast.T_VARIABLE_LENGTH_ARGUMENTS); tok != nil {
			arg.VariableLength = true
			break
		} else {
			break
		}
	}
	parser.expect(ast.T_PAREN_CLOSE)
	return args
}

func (parser *Parser) ParseIncludeStmt() ast.Stmt {
	var tok = parser.expect(ast.T_INCLUDE)
	var stm = ast.NewIncludeStmtWithToken(tok)

	var tok2 = parser.next()

	if tok2.Type == ast.T_IDENT {

		stm.MixinIdent = tok2

	} else if tok2.Type == ast.T_FUNCTION_NAME {

		stm.MixinIdent = tok2

		// TODO: revisit here later
		stm.ArgumentList = parser.ParseFunctionPrototype()

	} else {
		// TODO: report syntax error
		panic("Unexpected token after @include.")
	}

	var tok3 = parser.peek()
	if tok3.Type == ast.T_BRACE_OPEN {
		stm.ContentBlock = parser.ParseDeclBlock()
	}

	parser.expect(ast.T_SEMICOLON)
	return stm
}

func (parser *Parser) ParseFontFaceStmt() ast.Stmt {
	parser.expect(ast.T_FONT_FACE)
	block := parser.ParseDeclBlock()
	return &ast.FontFaceStmt{Block: block}
}

func (parser *Parser) ParseLogStmt() ast.Stmt {
	if directiveTok := parser.acceptAnyOf3(ast.T_ERROR, ast.T_WARN, ast.T_INFO); directiveTok != nil {
		var expr = parser.ParseExpr(false)
		parser.expect(ast.T_SEMICOLON)
		return &ast.LogStmt{
			Directive: directiveTok,
			Expr:      expr,
		}

	}
	panic(SyntaxError{
		Reason:      "Expecting @error, @warn, @info directive",
		ActualToken: parser.peek(),
	})
}

/*
@content directive is only allowed in mixin block
*/
func (parser *Parser) ParseContentStmt() ast.Stmt {
	var tok = parser.expect(ast.T_CONTENT)
	parser.expect(ast.T_SEMICOLON)
	return ast.NewContentStmtWithToken(tok)
}
