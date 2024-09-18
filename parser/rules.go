package parser

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/c9s/c6/ast"
	"github.com/c9s/c6/lexer"
)

var HttpUrlPattern = regexp.MustCompile("^https?://")
var AbsoluteUrlPattern = regexp.MustCompile("^[a-zA-Z]+?://")

func (parser *Parser) ReadFile(file string) error {
	f, err := ast.NewFile(parser.fsys, file)
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
	tokens, err := l.Run()

	if err != nil {
		return nil, err
	}

	parser.Tokens = tokens
	return parser.ParseStmts()
}

func (parser *Parser) ParseScss(code string) (*ast.StmtList, error) {
	l := lexer.NewLexerWithString(code)
	tokens, err := l.Run()

	if err != nil {
		return nil, err
	}

	parser.Tokens = tokens
	return parser.ParseStmts()
}

/*
ParseBlock method allows root level statements, which does not allow css properties.
*/
func (parser *Parser) ParseBlock() (*ast.Block, error) {
	debug("ParseBlock")
	if _, err := parser.expect(ast.T_BRACE_OPEN); err != nil {
		return nil, err
	}
	var block = ast.NewBlock()
	st, err := parser.ParseStmts()
	if err != nil {
		return nil, err
	}

	block.Stmts = st
	if _, err := parser.expect(ast.T_BRACE_CLOSE); err != nil {
		return nil, err
	}
	return block, nil
}

func (parser *Parser) ParseStmts() (*ast.StmtList, error) {
	var stmts = new(ast.StmtList)
	// stop at t_brace end
	for !parser.eof() {
		if stm, err := parser.ParseStmt(); err != nil {
			return nil, err
		} else if stm != nil {
			stmts.Append(stm)
		} else {
			break
		}
	}
	return stmts, nil
}

func (parser *Parser) ParseStmt() (ast.Stmt, error) {
	var token = parser.peek()

	if token == nil {
		return nil, nil
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
	case ast.T_AT_ROOT:
		return parser.ParseAtRootStmt()
	case ast.T_ERROR, ast.T_WARN, ast.T_DEBUG:
		return parser.ParseLogStmt()
	case ast.T_BRACKET_CLOSE:
		return nil, nil
	}

	if token.IsSelector() || token.IsSelectorCombinator() {
		return parser.ParseRuleSet()
	}

	// TODO: master just returns nil there, we'll do that too
	//return nil, fmt.Errorf("statement parse failed, unknown token [%s]", parser.peek())
	return nil, nil
}

func (parser *Parser) ParseIfStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_IF); err != nil {
		return nil, err
	}

	condition, err := parser.ParseCondition()
	if err != nil {
		return nil, err
	}

	if condition == nil {
		return nil, fmt.Errorf("if statement syntax error")
	}

	block, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}

	var stm = ast.NewIfStmt(condition, block)

	// TODO: OptimizeIfStmt(...)

	// If these is more else if statement
	var tok = parser.peek()
	for tok != nil && tok.Type == ast.T_ELSE_IF {
		parser.advance()

		condition, err := parser.ParseCondition()
		if err != nil {
			return nil, err
		}

		elseifblock, err := parser.ParseDeclBlock()
		if err != nil {
			return nil, err
		}

		var elseIfStm = ast.NewIfStmt(condition, elseifblock)
		stm.AppendElseIf(elseIfStm)
		tok = parser.peek()
	}

	tok = parser.peek()
	if tok != nil && tok.Type == ast.T_ELSE {
		parser.advance()

		// XXX: handle error here
		if elseBlock, err := parser.ParseDeclBlock(); err != nil {
			return nil, err
		} else if elseBlock != nil {
			stm.ElseBlock = elseBlock
		} else {
			return nil, SyntaxError{
				Reason:      "Expecting declaration block { ... }",
				ActualToken: parser.peek(),
				File:        parser.File,
			}
		}
	}

	return stm, nil
}

/*
The operator precedence is described here

@see http://introcs.cs.princeton.edu/java/11precedence/
*/
func (parser *Parser) ParseCondition() (ast.Expr, error) {
	debug("ParseCondition")

	// Boolean 'Not'
	if tok := parser.accept(ast.T_LOGICAL_NOT); tok != nil {
		logicexpr, err := parser.ParseLogicExpr()
		if err != nil {
			return nil, err
		}

		return ast.NewUnaryExpr(ast.NewOpWithToken(tok), logicexpr), nil
	}
	return parser.ParseLogicExpr()
}

func (parser *Parser) ParseLogicExpr() (ast.Expr, error) {
	debug("ParseLogicExpr")
	expr, err := parser.ParseLogicANDExpr()
	if err != nil {
		return nil, err
	}

	for tok := parser.accept(ast.T_LOGICAL_OR); tok != nil; tok = parser.accept(ast.T_LOGICAL_OR) {
		if subexpr, err := parser.ParseLogicANDExpr(); err != nil {
			return nil, err
		} else if subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
	}
	return expr, nil
}

func (parser *Parser) ParseLogicANDExpr() (ast.Expr, error) {
	debug("ParseLogicANDExpr")

	expr, err := parser.ParseComparisonExpr()
	if err != nil {
		return nil, err
	}

	for tok := parser.accept(ast.T_LOGICAL_AND); tok != nil; tok = parser.accept(ast.T_LOGICAL_AND) {
		if subexpr, err := parser.ParseComparisonExpr(); err != nil {
			return nil, err
		} else if subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
	}
	return expr, nil
}

func (parser *Parser) ParseComparisonExpr() (ast.Expr, error) {
	debug("ParseComparisonExpr")

	var expr ast.Expr
	var err error
	if parser.accept(ast.T_PAREN_OPEN) != nil {
		expr, err = parser.ParseLogicExpr()
		if err != nil {
			return nil, err
		}

		if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
			return nil, err
		}
	} else {
		expr, err = parser.ParseExpr(false)
		if err != nil {
			return nil, err
		}
	}

	var tok = parser.peek()
	for tok != nil && tok.IsComparisonOperator() {
		parser.advance()
		if subexpr, err := parser.ParseExpr(false); err != nil {
			return nil, err
		} else if subexpr != nil {
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
		tok = parser.peek()
	}
	return expr, nil
}

func (parser *Parser) ParseSimpleSelector(parentRuleSet *ast.RuleSet, pos int) (ast.Selector, error) {
	debug("ParseSimpleSelector")

	var tok = parser.next()
	if tok == nil {
		return nil, nil
	}

	switch tok.Type {

	case ast.T_TYPE_SELECTOR:

		return ast.NewTypeSelectorWithToken(tok), nil

	case ast.T_UNIVERSAL_SELECTOR:

		return ast.NewUniversalSelectorWithToken(tok), nil

	case ast.T_ID_SELECTOR:

		return ast.NewIdSelectorWithToken(tok), nil

	case ast.T_CLASS_SELECTOR:

		return ast.NewClassSelectorWithToken(tok), nil

	case ast.T_PARENT_SELECTOR:
		if pos > 0 {
			return nil, SyntaxError{
				Reason:      `"&" may only used at the beginning of a compound selector.`,
				ActualToken: parser.peek(),
				File:        parser.File,
			}
		}

		return ast.NewParentSelectorWithToken(tok), nil

	case ast.T_FUNCTIONAL_PSEUDO:

		var sel = ast.NewFunctionalPseudoSelectorWithToken(tok)
		if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
			return nil, err
		}

		var tok2 = parser.next()
		for tok2 != nil && tok2.Type != ast.T_PAREN_CLOSE {
			// TODO: parse pseudo expression
			tok2 = parser.next()
		}
		parser.backup()
		if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
			return nil, err
		}

		return sel, nil

	case ast.T_PSEUDO_SELECTOR:

		return ast.NewPseudoSelectorWithToken(tok), nil

	// Attribute selector parsing
	case ast.T_BRACKET_OPEN:
		attrName, err := parser.expect(ast.T_ATTRIBUTE_NAME)
		if err != nil {
			return nil, err
		}

		var tok2 = parser.next()
		if tok2.IsAttributeMatchOperator() {

			var tok3 = parser.next()
			var sel = ast.NewAttributeSelector(attrName, tok2, tok3)
			if _, err := parser.expect(ast.T_BRACKET_CLOSE); err != nil {
				return nil, err
			}

			return sel, nil

		} else if tok2.Type == ast.T_BRACKET_CLOSE {

			return ast.NewAttributeSelectorNameOnly(attrName), nil

		} else {

			return nil, SyntaxError{
				Reason:      "Unexpected token",
				ActualToken: tok2,
				File:        parser.File,
			}
		}
	}
	parser.backup()
	return nil, nil
}

func (parser *Parser) ParseCompoundSelector(parentRuleSet *ast.RuleSet) (*ast.CompoundSelector, error) {
	debug("ParseCompoundSelector")

	var sels = ast.NewCompoundSelector()
	idx := 0
	for {
		sel, err := parser.ParseSimpleSelector(parentRuleSet, idx)

		if err != nil {
			return nil, err
		}

		if sel != nil {
			sels.Append(sel)
			idx++
		} else {
			break
		}
	}
	if sels.Length() > 0 {
		return sels, nil
	}
	return nil, nil
}

func (parser *Parser) ParseComplexSelector(parentRuleSet *ast.RuleSet) (*ast.ComplexSelector, error) {
	debug("ParseComplexSelector")

	var complexSel = ast.NewComplexSelector()

	for {
		var tok = parser.peek()
		if tok == nil {
			return complexSel, nil
		}

		var comb ast.Combinator

		// peek the combinator token
		switch tok.Type {

		case ast.T_ADJACENT_SIBLING_COMBINATOR:

			comb = ast.NewAdjacentCombinatorWithToken(tok)
			parser.next()

		case ast.T_CHILD_COMBINATOR:

			comb = ast.NewChildCombinatorWithToken(tok)
			parser.next()

		case ast.T_DESCENDANT_COMBINATOR:

			comb = ast.NewDescendantCombinatorWithToken(tok)
			parser.next()

		case ast.T_GENERAL_SIBLING_COMBINATOR:

			comb = ast.NewGeneralSiblingCombinatorWithToken(tok)
			parser.next()

		default:
			if len(complexSel.ComplexSelectorItems) > 0 {
				return complexSel, nil
			}
		}

		if sel, err := parser.ParseCompoundSelector(parentRuleSet); err != nil {
			return nil, err
		} else if sel != nil || comb != nil {
			complexSel.AppendCompoundSelector(comb, sel)
		} else {
			// nothing found
			return nil, nil
		}
	}
}

func (parser *Parser) ParseSelectorList() (*ast.ComplexSelectorList, error) {
	debug("ParseSelectorList")

	var parentRuleSet = parser.GlobalContext.TopRuleSet()

	var complexSelectorList = &ast.ComplexSelectorList{}

	complexSelector, err := parser.ParseComplexSelector(parentRuleSet)

	if err != nil {
		return nil, err
	}

	if complexSelector != nil {
		complexSelectorList.Append(complexSelector)
	} else {
		return nil, nil
	}

	// if there is more comma
	for parser.accept(ast.T_COMMA) != nil {
		complexSelector, err := parser.ParseComplexSelector(parentRuleSet)

		if err != nil {
			return nil, err
		}

		if complexSelector != nil {
			complexSelectorList.Append(complexSelector)
		} else {
			break
		}
	}

	return complexSelectorList, nil
}

func (parser *Parser) ParseExtendStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_EXTEND); err != nil {
		return nil, err
	}
	var stm = ast.NewExtendStmt()
	selectors, err := parser.ParseSelectorList()
	if err != nil {
		return nil, err
	}

	stm.Selectors = selectors
	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}
	return stm, nil
}

func (parser *Parser) ParseRuleSet() (ast.Stmt, error) {
	var ruleset = ast.NewRuleSet()
	selectors, err := parser.ParseSelectorList()

	if err != nil {
		return nil, err
	}

	ruleset.Selectors = selectors

	parser.GlobalContext.PushRuleSet(ruleset)
	bl, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}

	ruleset.Block = bl
	parser.GlobalContext.PopRuleSet()

	return ruleset, nil
}

func (parser *Parser) ParseBoolean() ast.Expr {
	if tok := parser.acceptAnyOf2(ast.T_TRUE, ast.T_FALSE); tok != nil {
		return ast.NewBooleanWithToken(tok)
	}
	return nil
}

func (parser *Parser) ParseNumber() (ast.Expr, error) {
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
			return nil, err
		}
		if negative {
			i = -i
		}
		val = float64(i)

	} else if tok.Type == ast.T_FLOAT {

		f, err := strconv.ParseFloat(tok.Str, 64)
		if err != nil {
			return nil, err
		}
		if negative {
			f = -f
		}
		val = f

	} else {
		// unknown token
		parser.restore(pos)
		return nil, nil
	}

	if tok2.IsUnit() {
		// consume the unit token
		parser.next()
		return ast.NewNumber(val, ast.NewUnitWithToken(tok2), tok), nil
	}
	return ast.NewNumber(val, nil, tok), nil
}

/*
This parse method looks up the argument in the function declaration and convert
the keyword arguments into the argument list. this way, we can handle the function call
in a simple way - push/pop the stack.

@param fcall ast.FunctionCall The current parsing function call ast node
@return arguments []Expr
*/
func (parser *Parser) ParseKeywordArguments(fcall *ast.FunctionCall) (*ast.FunctionCallArguments, error) {
	// look up function declaration
	var item, ok = parser.GlobalContext.Functions.Get(fcall.Ident.Str)
	if !ok {
		return nil, fmt.Errorf("Undefined function %s", fcall.Ident.Str)
	}
	var fun = item.(*ast.Function)

	var arguments = ast.FunctionCallArguments{}
	for tok := parser.accept(ast.T_VARIABLE); tok != nil; tok = parser.accept(ast.T_VARIABLE) {
		var argdef = fun.ArgumentList.Lookup(tok.Str)
		if argdef == nil {
			return nil, fmt.Errorf("Undefined function argument: %s", tok.Str)
		}

		if _, err := parser.expect(ast.T_COLON); err != nil {
			return nil, err
		}

		argExpr, err := parser.ParseExpr(false)
		if err != nil {
			return nil, err
		}

		// TODO: this should be checked?
		parser.accept(ast.T_COMMA)

		var arg = ast.NewFunctionCallArgument(argExpr)
		arg.ArgumentDefineReference = argdef

		arguments = append(arguments, arg)
	}

	// sort arguments by ArgumentDefineReference
	arguments.Sort()
	return &arguments, nil
}

func (parser *Parser) ParseFunctionCall() (*ast.FunctionCall, error) {
	var identTok = parser.next()

	debug("ParseFunctionCall => next: %s", identTok)

	var fcall = ast.NewFunctionCallWithToken(identTok)

	if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
		return nil, err
	}

	var tok = parser.peek()
	var tok2 = parser.peekBy(2)
	if tok.Type == ast.T_VARIABLE && tok2.Type == ast.T_COLON {

		args, err := parser.ParseKeywordArguments(fcall)
		if err != nil {
			return nil, err
		}

		fcall.Arguments = *args

	} else {
		for tok.Type != ast.T_PAREN_CLOSE {

			if arg, err := parser.ParseFactor(); err != nil {
				return nil, err
			} else if arg != nil {
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
	if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
		return nil, err
	}

	return fcall, nil
}

func (parser *Parser) ParseIdent() (*ast.Ident, error) {
	var tok = parser.next()
	if tok.Type != ast.T_IDENT {
		return nil, fmt.Errorf("Invalid token for ident.")
	}
	return ast.NewIdentWithToken(tok), nil
}

/*
*
The ParseFactor must return an Expr interface compatible object
*/
func (parser *Parser) ParseFactor() (ast.Expr, error) {
	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_OPEN {

		if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
			return nil, err
		}

		expr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}

		if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
			return nil, err
		}
		return expr, nil

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	} else if tok.Type == ast.T_QQ_STRING {

		parser.advance()
		return ast.NewStringWithQuote('"', tok), nil

	} else if tok.Type == ast.T_Q_STRING {

		parser.advance()
		return ast.NewStringWithQuote('\'', tok), nil

	} else if tok.Type == ast.T_UNQUOTE_STRING {

		parser.advance()
		return ast.NewStringWithQuote(0, tok), nil

	} else if tok.Type == ast.T_TRUE {

		parser.advance()
		return ast.NewBooleanTrue(tok), nil

	} else if tok.Type == ast.T_FALSE {

		parser.advance()
		return ast.NewBooleanFalse(tok), nil

	} else if tok.Type == ast.T_NULL {

		parser.advance()
		return ast.NewNullWithToken(tok), nil

	} else if tok.Type == ast.T_FUNCTION_NAME {

		fcall, err := parser.ParseFunctionCall()
		if err != nil {
			return nil, err
		}
		return ast.Expr(fcall), nil

	} else if tok.Type == ast.T_VARIABLE {

		return parser.ParseVariable()

	} else if tok.Type == ast.T_IDENT {

		var tok2 = parser.peekBy(2)
		if tok2 != nil && tok2.Type == ast.T_PAREN_OPEN {
			return parser.ParseFunctionCall()
		}

		parser.advance()
		return ast.NewStringWithToken(tok), nil

	} else if tok.Type == ast.T_HEX_COLOR {

		parser.advance()
		return ast.NewHexColorFromToken(tok), nil

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		return parser.ParseNumber()

	}
	return nil, nil
}

func (parser *Parser) ParseTerm() (ast.Expr, error) {
	var pos = parser.Pos
	factor, err := parser.ParseFactor()
	if err != nil {
		return nil, err
	}

	if factor == nil {
		parser.restore(pos)
		return nil, nil
	}

	// see if the next token is '*' or '/'
	if tok := parser.acceptAnyOf2(ast.T_MUL, ast.T_DIV); tok != nil {
		term, err := parser.ParseTerm()

		if err != nil {
			return nil, err
		} else if term != nil {
			return ast.NewBinaryExpr(ast.NewOpWithToken(tok), factor, term, false), nil
		} else {
			return nil, SyntaxError{
				Reason:      "Expecting term after '*' or '/'",
				ActualToken: parser.peek(),
				File:        parser.File,
			}
		}
	}
	return factor, nil
}

/*
*

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};
*/
func (parser *Parser) ParseExpr(inParenthesis bool) (ast.Expr, error) {
	var pos = parser.Pos

	// plus or minus. This creates an unary expression that holds the later term.
	// this is for:  +3 or -4
	var expr ast.Expr
	var err error

	if tok := parser.acceptAnyOf2(ast.T_PLUS, ast.T_MINUS); tok != nil {
		if term, err := parser.ParseTerm(); err != nil {
			return nil, err
		} else if term != nil {
			expr = ast.NewUnaryExpr(ast.NewOpWithToken(tok), term)
		} else {
			parser.restore(pos)
			return nil, nil
		}
	} else {
		expr, err = parser.ParseTerm()
		if err != nil {
			return nil, err
		}
	}

	if expr == nil {
		debug("ParseExpr failed, got %+v, restoring to %d", expr, pos)
		parser.restore(pos)
		return nil, nil
	}

	var rightTok = parser.peek()
	for rightTok.Type == ast.T_PLUS || rightTok.Type == ast.T_MINUS || rightTok.Type == ast.T_LITERAL_CONCAT {
		// accept plus or minus
		parser.advance()

		if rightTerm, err := parser.ParseTerm(); err != nil {
			return nil, err
		} else if rightTerm != nil {
			// XXX: check parenthesis
			expr = ast.NewBinaryExpr(ast.NewOpWithToken(rightTok), expr, rightTerm, inParenthesis)
		} else {
			return nil, SyntaxError{
				Reason:      "Expecting term on the right side",
				ActualToken: parser.peek(),
				File:        parser.File,
			}
		}
		rightTok = parser.peek()
	}
	return expr, nil
}

func (parser *Parser) ParseMap() (ast.Expr, error) {
	var pos = parser.Pos
	var tok = parser.accept(ast.T_PAREN_OPEN)
	// since it's not started with '(', it's not map
	if tok == nil {
		parser.restore(pos)
		return nil, nil
	}

	var mapval = ast.NewMap()

	// TODO: check and report Map syntax error
	tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		keyExpr, err := parser.ParseExpr(false)
		if err != nil {
			return nil, err
		}
		if keyExpr == nil {
			parser.restore(pos)
			return nil, nil
		}

		if parser.accept(ast.T_COLON) == nil {
			parser.restore(pos)
			return nil, nil
		}

		valueExpr, err := parser.ParseExpr(false)
		if err != nil {
			return nil, err
		}
		if valueExpr == nil {
			parser.restore(pos)
			return nil, nil
		}

		// register the map value
		mapval.Set(keyExpr, valueExpr)
		parser.accept(ast.T_COMMA)
		tok = parser.peek()
	}
	if parser.accept(ast.T_PAREN_CLOSE) == nil {
		return nil, nil
	}
	return mapval, nil
}

func (parser *Parser) ParseString() (ast.Expr, error) {
	if tok := parser.accept(ast.T_QQ_STRING); tok != nil {

		return ast.NewStringWithQuote('"', tok), nil

	} else if tok := parser.accept(ast.T_Q_STRING); tok != nil {

		return ast.NewStringWithQuote('\'', tok), nil

	} else if tok := parser.accept(ast.T_UNQUOTE_STRING); tok != nil {

		return ast.NewStringWithQuote(0, tok), nil

	} else if tok := parser.accept(ast.T_IDENT); tok != nil {

		return ast.NewStringWithToken(tok), nil

	}

	var tok = parser.peek()
	if tok.Type == ast.T_INTERPOLATION_START {
		return parser.ParseInterp()
	}
	return nil, nil
}

func (parser *Parser) ParseInterp() (ast.Expr, error) {
	startTok, err := parser.expect(ast.T_INTERPOLATION_START)
	if err != nil {
		return nil, err
	}
	innerExpr, err := parser.ParseExpr(true)
	if err != nil {
		return nil, err
	}
	endTok, err := parser.expect(ast.T_INTERPOLATION_END)
	if err != nil {
		return nil, err
	}
	return ast.NewInterpolation(innerExpr, startTok, endTok), nil
}

func (parser *Parser) ParseValueStrict() (ast.Expr, error) {
	var pos = parser.Pos

	if tok := parser.accept(ast.T_PAREN_OPEN); tok != nil {
		if mapValue, err := parser.ParseMap(); err != nil {
			return nil, err
		} else if mapValue != nil {
			return mapValue, nil
		}
		parser.restore(pos)

		if listValue, err := parser.ParseList(); err != nil {
			return nil, err
		} else if listValue != nil {
			return listValue, nil
		}
		parser.restore(pos)
	}
	return parser.ParseExpr(false)
}

/*
Parse string literal expression (literal concat with interpolation)
*/
func (parser *Parser) ParseLiteralExpr() (ast.Expr, error) {
	if expr, err := parser.ParseExpr(false); err != nil {
		return nil, err
	} else if expr != nil {
		for tok := parser.accept(ast.T_LITERAL_CONCAT); tok != nil; tok = parser.accept(ast.T_LITERAL_CONCAT) {
			rightExpr, err := parser.ParseExpr(false)
			if err != nil {
				return nil, err
			} else if rightExpr == nil {
				return nil, SyntaxError{
					Reason:      "Expecting expression or ident after the literal concat operator.",
					ActualToken: parser.peek(),
					File:        parser.File,
				}
			}
			expr = ast.NewLiteralConcat(expr, rightExpr)
		}

		return expr, nil
	}
	return nil, nil
}

/*
*
Parse Value loosely

This parse method allows space/comma separated tokens of list.

To parse mixin argument or function argument, we only allow comma-separated list inside the parenthesis.

@param stopTokType ast.TokenType

The stop token is used from variable assignment expression,
We expect ';' semicolon at the end of expression to avoid the ambiguity of list, map and expression.
*/
func (parser *Parser) ParseValue(stopTokType ast.TokenType) (ast.Expr, error) {
	var pos = parser.Pos

	// try parse map
	debug("ParseMap")
	if mapValue, err := parser.ParseMap(); err != nil {
		return nil, err
	} else if mapValue != nil {
		var tok = parser.peek()
		if stopTokType == 0 || tok.Type == stopTokType {
			debug("OK Map Meet Stop Token")
			return mapValue, nil
		}
	}
	debug("Map parse failed, restoring to %d", pos)
	parser.restore(pos)

	debug("Trying List")
	if listValue, err := parser.ParseList(); err != nil {
		return nil, err
	} else if listValue != nil {
		var tok = parser.peek()
		if stopTokType == 0 || tok.Type == stopTokType {
			debug("OK List: %+v", listValue)
			return listValue, nil
		}
	}

	debug("List parse failed, restoring to %d", pos)
	parser.restore(pos)

	debug("ParseLiteralExpr trying", pos)
	return parser.ParseLiteralExpr()
}

func (parser *Parser) ParseList() (ast.Expr, error) {
	debug("ParseList at %d", parser.Pos)
	var pos = parser.Pos
	if list, err := parser.ParseCommaSepList(); err != nil {
		return nil, err
	} else if list != nil {
		return list, nil
	}
	parser.restore(pos)
	return nil, nil
}

func (parser *Parser) ParseCommaSepList() (ast.Expr, error) {
	debug("ParseCommaSepList at %d", parser.Pos)
	var list = ast.NewCommaSepList()

	var tok = parser.peek()
	for tok.Type != ast.T_COMMA && tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {

		// when the syntax start with a '(', it could be a list or map.
		if tok.Type == ast.T_PAREN_OPEN {

			parser.next()
			if sublist, err := parser.ParseCommaSepList(); err != nil {
				return nil, err
			} else if sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			}
			// allow empty list here
			if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
				return nil, err
			}

		} else {
			sublist, err := parser.ParseSpaceSepList()
			if err != nil {
				return nil, err
			}
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

		return nil, nil

	} else if list.Len() == 1 {

		return list.Exprs[0], nil

	}
	return list, nil
}

func (parser *Parser) ParseVariable() (*ast.Variable, error) {
	if tok := parser.accept(ast.T_VARIABLE); tok != nil {
		return ast.NewVariableWithToken(tok), nil
	}
	return nil, nil
}

func (parser *Parser) ParseAssignStmt() (ast.Stmt, error) {
	variable, err := parser.ParseVariable()
	if err != nil {
		return nil, err
	}

	// skip ":", T_COLON token
	if _, err := parser.expect(ast.T_COLON); err != nil {
		return nil, err
	}

	// Expecting semicolon at the end of the statement
	valExpr, err := parser.ParseValue(ast.T_SEMICOLON)
	if err != nil {
		return nil, err
	} else if valExpr == nil {
		return nil, SyntaxError{
			Reason:      "Expecting value after variable assignment.",
			ActualToken: parser.peek(),
			File:        parser.File,
		}
	}

	var stm = ast.NewAssignStmt(variable, valExpr)
	parser.ParseFlags(stm)
	parser.accept(ast.T_SEMICOLON)
	return stm, nil
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

func (parser *Parser) ParseSpaceSepList() (ast.Expr, error) {
	debug("ParseSpaceSepList at %d", parser.Pos)

	var list = ast.NewSpaceSepList()
	list.Separator = " "

	if tok := parser.accept(ast.T_PAREN_OPEN); tok != nil {
		if sublist, err := parser.ParseCommaSepList(); err != nil {
			return nil, err
		} else if sublist != nil {
			list.Append(sublist)
		}
		if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
			return nil, err
		}
	}

	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {
		subexpr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}
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
		return nil, nil
	} else if list.Len() == 1 {
		return list.Exprs[0], nil
	} else if list.Len() > 1 {
		return list, nil
	}
	return nil, nil
}

/*
*
We treat the property value section as a list value, which is separated by ',' or ' '
*/
func (parser *Parser) ParsePropertyValue(parentRuleSet *ast.RuleSet, property *ast.Property) (*ast.List, error) {
	debug("ParsePropertyValue")
	// var tok = parser.peek()
	var list = ast.NewSpaceSepList()

	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {
		sublist, err := parser.ParseList()
		if err != nil {
			return nil, err
		}
		if sublist != nil {
			list.Append(sublist)
			debug("ParsePropertyValue list: %+v", list)
		} else {
			break
		}
		tok = parser.peek()
	}

	return list, nil
}

func (parser *Parser) ParsePropertyName() (ast.Expr, error) {
	ident, err := parser.ParsePropertyNameToken()
	if err != nil {
		return nil, err
	}
	if ident == nil {
		return nil, nil
	}

	var tok = parser.peek()
	for tok.Type == ast.T_LITERAL_CONCAT {
		parser.next()
		_, err := parser.ParsePropertyNameToken()
		if err != nil {
			return nil, err
		}
		tok = parser.peek()
	}
	if _, err := parser.expect(ast.T_COLON); err != nil {
		return nil, err
	}
	return ident, nil // TODO: new literal concat ast
}

func (parser *Parser) ParsePropertyNameToken() (ast.Expr, error) {
	if tok := parser.accept(ast.T_PROPERTY_NAME_TOKEN); tok != nil {
		return ast.NewIdentWithToken(tok), nil
	}

	var tok = parser.peek()
	if tok.Type == ast.T_INTERPOLATION_START {
		return parser.ParseInterpolation()
	}
	return nil, nil
}

func (parser *Parser) ParseInterpolation() (ast.Expr, error) {
	debug("ParseInterpolation")
	var startToken *ast.Token
	if startToken = parser.accept(ast.T_INTERPOLATION_START); startToken == nil {
		return nil, nil
	}
	expr, err := parser.ParseExpr(true)
	if err != nil {
		return nil, err
	}
	endToken, err := parser.expect(ast.T_INTERPOLATION_END)
	if err != nil {
		return nil, err
	}
	return ast.NewInterpolation(expr, startToken, endToken), nil
}

func (parser *Parser) ParseDeclaration() ast.Stmt {
	return nil
}

func (parser *Parser) ParseDeclBlock() (*ast.DeclBlock, error) {
	var parentRuleSet = parser.GlobalContext.TopRuleSet()
	var declBlock = ast.NewDeclBlock(parentRuleSet)

	if _, err := parser.expect(ast.T_BRACE_OPEN); err != nil {
		return nil, err
	}

	var tok = parser.peek()
	for tok != nil && tok.Type != ast.T_BRACE_CLOSE {
		if propertyName, err := parser.ParsePropertyName(); err != nil {
			return nil, err
		} else if propertyName != nil {
			var property = ast.NewProperty(tok)

			if valueList, err := parser.ParsePropertyValue(parentRuleSet, property); err != nil {
				return nil, err
			} else if valueList != nil {
				for _, v := range valueList.Exprs {
					property.AppendValue(v)
				}
			}
			declBlock.Append(property)

			var tok2 = parser.peek()

			// if nested property found
			if tok2.Type == ast.T_BRACE_OPEN {
				// TODO: merge them back to current block
				_, err := parser.ParseDeclBlock()

				if err != nil {
					return nil, err
				}

				// TODO: why do we parse nested block if we don't do
				// a thing with it?
				//return nil, fmt.Errorf("TODO: nested declaration block not implemented")
			}

			if parser.accept(ast.T_SEMICOLON) == nil {
				if tok3 := parser.peek(); tok3.Type == ast.T_BRACE_CLOSE {
					// normal break
					break
				} else {
					return nil, fmt.Errorf("missing semicolon after the property value.")
				}
			}

		} else if stm, err := parser.ParseStmt(); err != nil {
			return nil, err
		} else if stm != nil {
			declBlock.Append(stm)
		} else {
			return nil, fmt.Errorf("Parse failed at token %s", tok)
		}
		tok = parser.peek()

	}
	if _, err := parser.expect(ast.T_BRACE_CLOSE); err != nil {
		return nil, err
	}
	return declBlock, nil
}

func (parser *Parser) ParseCharsetStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_CHARSET); err != nil {
		return nil, err
	}
	var tok = parser.next()
	var stm = ast.NewCharsetStmtWithToken(tok)
	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}
	return stm, nil
}

/*
Media Query Syntax:
https://developer.mozilla.org/en-US/docs/Web/Guide/CSS/Media_queries
*/
func (parser *Parser) ParseMediaQueryStmt() (ast.Stmt, error) {
	// expect the '@media' token
	var stm = ast.NewMediaQueryStmt()
	if _, err := parser.expect(ast.T_MEDIA); err != nil {
		return nil, err
	}

	if list, err := parser.ParseMediaQueryList(); err != nil {
		return nil, err
	} else if list != nil {
		stm.MediaQueryList = list
	}
	bl, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}
	stm.Block = bl
	return stm, nil
}

func (parser *Parser) ParseMediaQueryList() (*ast.MediaQueryList, error) {
	query, err := parser.ParseMediaQuery()
	if err != nil {
		return nil, err
	}

	if query == nil {
		return nil, nil
	}

	var queries = &ast.MediaQueryList{List: []*ast.MediaQuery{query}}
	for parser.accept(ast.T_COMMA) != nil {
		if query, err := parser.ParseMediaQuery(); err != nil {
			return nil, err
		} else if query != nil {
			queries.Append(query)
		}
	}
	return queries, nil
}

/*
This method parses media type first, then expecting more that on media
expressions.

media_query: [[only | not]? <media_type> [ and <expression> ]*]

	| <expression> [ and <expression> ]*

expression: ( <media_feature> [: <value>]? )

Specification: http://dev.w3.org/csswg/mediaqueries-3
*/
func (parser *Parser) ParseMediaQuery() (*ast.MediaQuery, error) {

	// the leading media type is optional
	mediaType, err := parser.ParseMediaType()
	if err != nil {
		return nil, err
	}
	if mediaType != nil {
		// Check if there is an expression after the media type.
		var tok = parser.peek()
		if tok.Type != ast.T_LOGICAL_AND {
			return ast.NewMediaQuery(mediaType, nil), nil
		}
		parser.advance() // skip the and operator token
	}

	// parse the media expression after the media type.
	mediaExpr, err := parser.ParseMediaQueryExpr()
	if err != nil {
		return nil, err
	}
	if mediaExpr == nil {
		if mediaType == nil {
			return nil, nil
		}
		return ast.NewMediaQuery(mediaType, mediaExpr), nil
	}

	// @media query only allows AND operator here..
	for tok := parser.accept(ast.T_LOGICAL_AND); tok != nil; tok = parser.accept(ast.T_LOGICAL_AND) {
		// parse another mediq query expression
		expr2, err := parser.ParseMediaQueryExpr()
		if err != nil {
			return nil, err
		}
		mediaExpr = ast.NewBinaryExpr(ast.NewOpWithToken(tok), mediaExpr, expr2, false)
	}
	return ast.NewMediaQuery(mediaType, mediaExpr), nil
}

/*
ParseMediaType returns Ident Node or UnaryExpr as ast.Expr
*/
func (parser *Parser) ParseMediaType() (*ast.MediaType, error) {
	if tok := parser.acceptAnyOf2(ast.T_LOGICAL_NOT, ast.T_ONLY); tok != nil {
		mediaType, err := parser.expect(ast.T_IDENT)
		if err != nil {
			return nil, err
		}

		return ast.NewMediaType(ast.NewUnaryExpr(ast.NewOpWithToken(tok), mediaType)), nil
	}

	var tok = parser.peek()
	if tok.Type == ast.T_PAREN_OPEN {
		// the begining of the media expression
		return nil, nil
	}

	expr, err := parser.ParseExpr(false)
	if err != nil {
		return nil, err
	}

	if expr != nil {
		return ast.NewMediaType(expr), nil
	}

	// parse media type fail
	return nil, nil
}

/*
An media query expression must start with a '(' and ends with ')'
*/
func (parser *Parser) ParseMediaQueryExpr() (ast.Expr, error) {

	// it's not an media query expression
	if openTok := parser.accept(ast.T_PAREN_OPEN); openTok != nil {
		featureExpr, err := parser.ParseExpr(false)
		if err != nil {
			return nil, err
		}
		var feature = ast.NewMediaFeature(featureExpr, nil)

		// if the next token is a colon, then we expect a feature value
		// after the colon.
		if tok := parser.accept(ast.T_COLON); tok != nil {
			val, err := parser.ParseExpr(false)
			if err != nil {
				return nil, err
			}
			feature.Value = val
		}
		closeTok, err := parser.expect(ast.T_PAREN_CLOSE)

		if err != nil {
			return nil, err
		}
		feature.Open = openTok
		feature.Close = closeTok
		return feature, nil
	}
	return nil, nil
}

func (parser *Parser) ParseWhileStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_WHILE); err != nil {
		return nil, err
	}

	condition, err := parser.ParseCondition()
	if err != nil {
		return nil, err
	}
	block, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}
	return ast.NewWhileStmt(condition, block), nil
}

/*
Parse the SASS @for statement.

	@for $var from <start> to <end> {  }

	@for $var from <start> through <end> {  }

@see http://sass-lang.com/documentation/file.SASS_REFERENCE.html#_10
*/
func (parser *Parser) ParseForStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_FOR); err != nil {
		return nil, err
	}

	// get the variable token
	variable, err := parser.ParseVariable()
	if err != nil {
		return nil, err
	}
	var stm = ast.NewForStmt(variable)

	if parser.accept(ast.T_FOR_FROM) != nil {

		fromExpr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}

		stm.From = fromExpr

		// "through" or "to"
		var tok = parser.next()

		if tok.Type != ast.T_FOR_THROUGH && tok.Type != ast.T_FOR_TO {
			return nil, SyntaxError{
				Reason:      "Expecting 'through' or 'to' of range syntax.",
				ActualToken: tok,
				File:        parser.File,
			}
		}

		endExpr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}

		if tok.Type == ast.T_FOR_THROUGH {

			stm.To = endExpr
			stm.Inclusive = true

		} else if tok.Type == ast.T_FOR_TO {

			stm.To = endExpr

		}

	} else if parser.accept(ast.T_FOR_IN) != nil {

		fromExpr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}
		stm.From = fromExpr

		if _, err := parser.expect(ast.T_RANGE); err != nil {
			return nil, err
		}

		endExpr, err := parser.ParseExpr(true)
		if err != nil {
			return nil, err
		}
		stm.To = endExpr
	}

	if b, err := parser.ParseDeclBlock(); err != nil {
		return nil, err
	} else if b != nil {
		stm.Block = b
	} else {
		return nil, fmt.Errorf("The @for statement expecting block after the range syntax")
	}
	return stm, nil
}

/*
The @import syntax is described here:

@see CSS2.1 http://www.w3.org/TR/CSS2/cascade.html#at-import

@see https://developer.mozilla.org/en-US/docs/Web/CSS/@import
*/
func (parser *Parser) ParseImportStmt() (ast.Stmt, error) {
	// skip the ast.T_IMPORT token
	if _, err := parser.expect(ast.T_IMPORT); err != nil {
		return nil, err
	}

	// Create the import statement node
	var stm = ast.NewImportStmt()

	var tok = parser.peek()

	// if it's url(..)
	if tok.Type == ast.T_FUNCTION_NAME {
		parser.advance()
		if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
			return nil, err
		}

		var urlTok = parser.acceptAnyOf3(ast.T_QQ_STRING, ast.T_Q_STRING, ast.T_UNQUOTE_STRING)
		if urlTok == nil {
			return nil, SyntaxError{
				Reason:      "Expecting url string in the url() function expression",
				ActualToken: parser.peek(),
				File:        parser.File,
			}
		}

		if HttpUrlPattern.MatchString(urlTok.Str) {
			stm.Url = ast.AbsoluteUrl(urlTok.Str)
		} else {
			stm.Url = ast.RelativeUrl(urlTok.Str)
		}

		if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
			return nil, err
		}

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
				return nil, fmt.Errorf("Unknown scss file to detect import path.")
			}

			var importPath = tok.Str
			fi, err := fs.Stat(parser.fsys, importPath)
			if err != nil {
				return nil, err
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

			if _, err := fs.Stat(parser.fsys, importPath); err != nil {
				return nil, err
			}
			stm.Url = ast.ScssImportUrl(importPath)

			if _, ok := parser.GlobalContext.ImportedPath[importPath]; !ok {
				// Set imported path to true
				parser.GlobalContext.ImportedPath[importPath] = true

				// parse the imported file using the same context
				var subparser = NewParser(parser.GlobalContext)
				var stmts, err = subparser.ParseScssFile(importPath)
				if err != nil {
					return nil, err
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
		return nil, fmt.Errorf("Unexpected token: %s", tok)
	}

	if mediaQueryList, err := parser.ParseMediaQueryList(); err != nil {
		return nil, err
	} else if mediaQueryList != nil {
		stm.MediaQueryList = mediaQueryList.List
	}

	// must be ast.T_SEMICOLON at the end
	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}
	return stm, nil
}

func (parser *Parser) ParseReturnStmt() (ast.Stmt, error) {
	returnTok, err := parser.expect(ast.T_RETURN)
	if err != nil {
		return nil, err
	}
	valueExpr, err := parser.ParseExpr(true)
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}

	return ast.NewReturnStmtWithToken(returnTok, valueExpr), nil
}

func (parser *Parser) ParseFunctionDeclaration() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_FUNCTION); err != nil {
		return nil, err
	}

	identTok, err := parser.expect(ast.T_FUNCTION_NAME)
	if err != nil {
		return nil, err
	}

	args, err := parser.ParseFunctionPrototype()
	if err != nil {
		return nil, err
	}

	var fun = ast.NewFunctionWithToken(identTok)
	fun.ArgumentList = args

	if bl, err := parser.ParseBlock(); err != nil {
		return nil, err
	} else {
		fun.Block = bl
	}
	parser.GlobalContext.Functions.Set(identTok.Str, fun)
	return fun, nil
}

func (parser *Parser) ParseMixinStmt() (ast.Stmt, error) {
	mixinTok, err := parser.expect(ast.T_MIXIN)
	if err != nil {
		return nil, err
	}

	var stm = ast.NewMixinStmtWithToken(mixinTok)

	var tok = parser.next()

	// Mixin without parameters
	if tok.Type == ast.T_IDENT {

		stm.Ident = tok

	} else if tok.Type == ast.T_FUNCTION_NAME {

		stm.Ident = tok
		l, err := parser.ParseFunctionPrototype()
		if err != nil {
			return nil, err
		}

		stm.ArgumentList = l

	} else {
		return nil, fmt.Errorf("Syntax error")
	}

	if b, err := parser.ParseDeclBlock(); err != nil {
		return nil, err
	} else {
		stm.Block = b
	}

	parser.GlobalContext.Mixins.Set(stm.Ident.Str, stm)
	return stm, nil
}

func (parser *Parser) ParseFunctionPrototypeArgument() (*ast.Argument, error) {
	debug("ParseFunctionPrototypeArgument")

	var varTok *ast.Token
	if varTok = parser.accept(ast.T_VARIABLE); varTok == nil {
		return nil, nil
	}

	if arg := ast.NewArgumentWithToken(varTok); arg != nil {
		if parser.accept(ast.T_COLON) != nil {
			v, err := parser.ParseValueStrict()
			if err != nil {
				return nil, err
			}
			arg.DefaultValue = v
		}
		return arg, nil
	}
	return nil, nil
}

func (parser *Parser) ParseFunctionPrototype() (*ast.ArgumentList, error) {
	debug("ParseFunctionPrototype")

	var args = ast.NewArgumentList()

	if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
		return nil, err
	}

	var tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var arg *ast.Argument
		var err error
		if arg, err = parser.ParseFunctionPrototypeArgument(); err != nil {
			return nil, err
		} else if arg != nil {
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
	if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
		return nil, err
	}

	return args, nil
}

func (parser *Parser) ParseFunctionCallArgument() (*ast.Argument, error) {
	debug("ParseFunctionCallArgument")

	var varTok *ast.Token
	if varTok = parser.accept(ast.T_VARIABLE); varTok == nil {
		return nil, nil
	}

	var arg = ast.NewArgumentWithToken(varTok)

	if parser.accept(ast.T_COLON) != nil {
		val, err := parser.ParseValueStrict()
		if err != nil {
			return nil, err
		}

		arg.DefaultValue = val
	}
	return arg, nil
}

func (parser *Parser) ParseFunctionCallArguments() (*ast.ArgumentList, error) {
	debug("ParseFunctionCallArguments")

	var args = ast.NewArgumentList()

	if _, err := parser.expect(ast.T_PAREN_OPEN); err != nil {
		return nil, err
	}
	var tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var arg *ast.Argument
		var err error
		if arg, err = parser.ParseFunctionCallArgument(); err != nil {
			return nil, err
		} else if arg != nil {
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
	if _, err := parser.expect(ast.T_PAREN_CLOSE); err != nil {
		return nil, err
	}
	return args, nil
}

func (parser *Parser) ParseIncludeStmt() (ast.Stmt, error) {
	tok, err := parser.expect(ast.T_INCLUDE)
	if err != nil {
		return nil, err
	}

	var stm = ast.NewIncludeStmtWithToken(tok)

	var tok2 = parser.next()

	if tok2.Type == ast.T_IDENT {

		stm.MixinIdent = tok2

	} else if tok2.Type == ast.T_FUNCTION_NAME {

		stm.MixinIdent = tok2

		if al, err := parser.ParseFunctionPrototype(); err != nil {
			return nil, err
		} else {
			stm.ArgumentList = al
		}

	} else {
		return nil, fmt.Errorf("Unexpected token after @include.")
	}

	var tok3 = parser.peek()
	if tok3.Type == ast.T_BRACE_OPEN {
		if bl, err := parser.ParseDeclBlock(); err != nil {
			return nil, err
		} else {
			stm.ContentBlock = bl
		}
	}

	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}

	return stm, nil
}

func (parser *Parser) ParseFontFaceStmt() (ast.Stmt, error) {
	if _, err := parser.expect(ast.T_FONT_FACE); err != nil {
		return nil, err
	}
	block, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FontFaceStmt{Block: block}, nil
}

func (parser *Parser) ParseLogStmt() (ast.Stmt, error) {
	if directiveTok := parser.acceptAnyOf3(ast.T_ERROR, ast.T_WARN, ast.T_DEBUG); directiveTok != nil {
		var ll ast.LogLevel

		switch directiveTok.Type {
		case ast.T_DEBUG:
			ll = ast.LogLevelDebug
		case ast.T_WARN:
			ll = ast.LogLevelWarn
		default:
			ll = ast.LogLevelError
		}

		expr, err := parser.ParseString()
		if err != nil {
			return nil, err
		}

		if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
			return nil, err
		}

		return ast.NewLogStmt(expr, directiveTok, ll), nil

	}
	return nil, SyntaxError{
		Reason:      "Expecting @error, @warn, @debug directive",
		ActualToken: parser.peek(),
	}
}

/*
@content directive is only allowed in mixin block
*/
func (parser *Parser) ParseContentStmt() (ast.Stmt, error) {
	tok, err := parser.expect(ast.T_CONTENT)
	if err != nil {
		return nil, err
	}

	if _, err := parser.expect(ast.T_SEMICOLON); err != nil {
		return nil, err
	}

	return ast.NewContentStmtWithToken(tok), nil
}

/*
@content directive is only allowed in mixin block
*/
func (parser *Parser) ParseAtRootStmt() (ast.Stmt, error) {
	tok, err := parser.expect(ast.T_AT_ROOT)
	if err != nil {
		return nil, err
	}

	stm := ast.NewAtRootStmtWithToken(tok)

	tok = parser.peek()

	if tok.IsSelector() {
		sel, err := parser.ParseComplexSelector(nil)

		if err != nil {
			return nil, err
		}

		stm.Selector = sel
	} else if tok.Type == ast.T_PAREN_OPEN {
		return nil, SyntaxError{
			Reason:      "@at-root does not support expressions like (without: media) yet",
			ActualToken: parser.peek(),
		}
	}

	bl, err := parser.ParseDeclBlock()
	if err != nil {
		return nil, err
	}
	stm.Block = bl

	return stm, nil
}
