package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "fmt"
import "strconv"
import "c6/ast"
import "c6/runtime"

func (parser *Parser) ParseBlock() *ast.Block {
	debug("ParseBlock")
	parser.expect(ast.T_BRACE_OPEN)
	var block = ast.NewBlock()
	block.Statements = parser.ParseStatements()
	parser.expect(ast.T_BRACE_CLOSE)
	return block
}

func (parser *Parser) ParseStatements() []ast.Statement {
	var stmts []ast.Statement = []ast.Statement{}
	// stop at t_brace end
	for !parser.eof() {
		if stmt := parser.ParseStatement(); stmt != nil {
			stmts = append(stmts, stmt)
		} else {
			break
		}
	}
	return stmts
}

func (parser *Parser) ParseStatement() ast.Statement {
	var token = parser.peek()

	if token == nil {
		return nil
	}

	switch token.Type {
	case ast.T_IMPORT:
		return parser.ParseImportStatement()
	case ast.T_CHARSET:
		return parser.ParseCharsetStatement()
	case ast.T_MEDIA:
		return parser.ParseMediaQueryStatement()
	case ast.T_MIXIN:
		return parser.ParseMixinStatement()
	case ast.T_FUNCTION:
		return parser.ParseFunction()
	case ast.T_INCLUDE:
		return parser.ParseIncludeStatement()
	case ast.T_VARIABLE:
		return parser.ParseVariableAssignment()
	case ast.T_RETURN:
		return parser.ParseReturnStatement()
	case ast.T_IF:
		return parser.ParseIfStatement()
	case ast.T_FOR:
		return parser.ParseForStatement()
	case ast.T_CONTENT:
		return parser.ParseContentStatement()
	}

	if token.IsSelector() {

		return parser.ParseRuleSet()

	} else {
		// panic(fmt.Errorf("parse failed, unknown token", parser.peek()))
	}
	return nil
}

func (parser *Parser) ParseIfStatement() ast.Statement {
	parser.expect(ast.T_IF)
	condition := parser.ParseCondition()
	if condition == nil {
		panic("if statement syntax error")
	}

	var block = parser.ParseBlock()
	var stm = ast.NewIfStatement(condition, block)

	// TODO: runtime.OptimizeIfStatement(...)

	// If these is more else if statement
	var tok = parser.peek()
	for tok != nil && tok.Type == ast.T_ELSE_IF {
		parser.next()

		// XXX: handle error here
		var condition = parser.ParseCondition()
		var elseifblock = parser.ParseBlock()
		var elseIfStm = ast.NewIfStatement(condition, elseifblock)
		stm.AppendElseIf(elseIfStm)
		tok = parser.peek()
	}

	tok = parser.peek()
	if tok != nil && tok.Type == ast.T_ELSE {
		parser.next()

		// XXX: handle error here
		var elseBlock = parser.ParseBlock()
		stm.ElseBlock = elseBlock
	}

	return stm
}

/*
The operator precedence is described here

@see http://introcs.cs.princeton.edu/java/11precedence/
*/
func (parser *Parser) ParseCondition() ast.Expression {
	debug("ParseCondition")

	// Boolean 'Not'
	var tok = parser.peek()
	if tok.Type == ast.T_LOGICAL_NOT {
		var logicexpr = parser.ParseLogicExpression()
		return ast.NewUnaryExpression(ast.NewOpWithToken(tok), logicexpr)
	}
	return parser.ParseLogicExpression()
}

func (parser *Parser) ParseLogicExpression() ast.Expression {
	debug("ParseLogicExpression")
	var expr = parser.ParseLogicANDExpression()
	var tok = parser.peek()
	for tok != nil && tok.Type == ast.T_LOGICAL_OR {
		parser.next()
		if subexpr := parser.ParseLogicANDExpression(); subexpr != nil {
			expr = ast.NewBinaryExpression(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
		tok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseLogicANDExpression() ast.Expression {
	debug("ParseLogicANDExpression")
	var expr = parser.ParseComparisonExpression()
	var tok = parser.peek()
	for tok != nil && tok.Type == ast.T_LOGICAL_AND {
		parser.next()
		if subexpr := parser.ParseComparisonExpression(); subexpr != nil {
			expr = ast.NewBinaryExpression(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
		tok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseComparisonExpression() ast.Expression {
	debug("ParseComparisonExpression")

	var expr ast.Expression
	var tok = parser.peek()
	if tok.Type == ast.T_PAREN_OPEN {
		parser.accept(ast.T_PAREN_OPEN)
		expr = parser.ParseLogicExpression()
		parser.expect(ast.T_PAREN_CLOSE)
	} else {
		expr = parser.ParseExpression(false)
	}

	tok = parser.peek()
	for tok != nil && tok.IsComparisonOperator() {
		parser.next()
		if subexpr := parser.ParseExpression(false); subexpr != nil {
			expr = ast.NewBinaryExpression(ast.NewOpWithToken(tok), expr, subexpr, false)
		}
		tok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseSelectors() *ast.SelectorList {
	var parentRuleSet = parser.Context.TopRuleSet()
	var sels = ast.NewSelectorList()

	var tok = parser.next()

	// TODO: Strict checking with selector combinator
	for tok.IsSelector() || tok.IsSelectorCombinator() {
		switch tok.Type {

		case ast.T_TYPE_SELECTOR:

			var sel = ast.NewTypeSelectorWithToken(tok)
			sels.Append(sel)

		case ast.T_UNIVERSAL_SELECTOR:

			var sel = ast.NewUniversalSelectorWithToken(tok)
			sels.Append(sel)

		case ast.T_ID_SELECTOR:

			var sel = ast.NewIdSelectorWithToken(tok)
			sels.Append(sel)

		case ast.T_CLASS_SELECTOR:

			var sel = ast.NewClassSelectorWithToken(tok)
			sels.Append(sel)

		case ast.T_PARENT_SELECTOR:

			var sel = ast.NewParentSelectorWithToken(parentRuleSet, tok)
			sels.Append(sel)

		case ast.T_PSEUDO_SELECTOR:

			// pseudo selector
			var sel = ast.NewPseudoSelectorWithToken(tok)
			if nextTok := parser.peek(); nextTok.Type == ast.T_LANG_CODE {
				sel.C = nextTok.Str
			}
			sels.Append(sel)

		case ast.T_ADJACENT_SIBLING_COMBINATOR:

			var sel = ast.NewAdjacentCombinatorWithToken(tok)
			sels.Append(sel)

		case ast.T_CHILD_COMBINATOR:

			var sel = ast.NewChildCombinatorWithToken(tok)
			sels.Append(sel)

		case ast.T_DESCENDANT_COMBINATOR:

			var sel = ast.NewDescendantCombinatorWithToken(tok)
			sels.Append(sel)

		case ast.T_COMMA:

			var sel = ast.NewGroupCombinatorWithToken(tok)
			sels.Append(sel)

		case ast.T_BRACKET_OPEN:
			var attrName = parser.expect(ast.T_ATTRIBUTE_NAME)
			var tok2 = parser.next()
			if tok2.IsAttributeMatchOperator() {
				var tok3 = parser.next()
				var sel = ast.NewAttributeSelector(attrName, tok2, tok3)
				sels.Append(sel)
				parser.expect(ast.T_BRACKET_CLOSE)

			} else if tok2.Type == ast.T_BRACKET_CLOSE {
				var sel = ast.NewAttributeSelectorNameOnly(attrName)
				sels.Append(sel)

			} else {
				panic(fmt.Errorf("Unexpected token type: %s", tok2))
			}

		default:
			panic(fmt.Errorf("Unexpected selector token: %+v", tok))
		}
		tok = parser.next()
	}
	parser.backup()
	return sels
}

func (parser *Parser) ParseRuleSet() ast.Statement {
	var ruleset = ast.NewRuleSet()

	ruleset.Selectors = parser.ParseSelectors()

	parser.Context.PushRuleSet(ruleset)

	ruleset.Block = parser.ParseDeclarationBlock()

	parser.Context.PopRuleSet()
	return ruleset
}

func (parser *Parser) ParseBoolean() ast.Expression {
	var tok = parser.peek()

	if tok.Type == ast.T_TRUE {
		parser.next()
		return ast.NewBooleanTrue(tok)
	} else if tok.Type == ast.T_FALSE {
		parser.next()
		return ast.NewBooleanFalse(tok)
	}
	return nil
}

func (parser *Parser) ParseNumber() ast.Expression {
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

func (parser *Parser) ParseFunctionCall() *ast.FunctionCall {
	var identTok = parser.next()

	debug("ParseFunctionCall => next: %s", identTok)

	var fcall = ast.NewFunctionCallWithToken(identTok)

	parser.expect(ast.T_PAREN_OPEN)

	var tok = parser.peek()
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
	parser.expect(ast.T_PAREN_CLOSE)
	return fcall
}

func (parser *Parser) ParseIdent() *ast.Ident {
	var tok = parser.next()
	debug("ReduceIndent => next: %s", tok)
	if tok.Type != ast.T_IDENT {
		panic("Invalid token for ident.")
	}
	return ast.NewIdentWithToken(tok)
}

/**
The ParseFactor must return an Expression interface compatible object
*/
func (parser *Parser) ParseFactor() ast.Expression {
	debug("ParseFactor at %d", parser.Pos)
	var tok = parser.peek()
	debug("ParseFactor => peek: %s", tok)

	if tok.Type == ast.T_PAREN_OPEN {
		parser.expect(ast.T_PAREN_OPEN)
		var expr = parser.ParseExpression(true)
		parser.expect(ast.T_PAREN_CLOSE)
		return expr

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	} else if tok.Type == ast.T_QQ_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('"', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_Q_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('\'', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_TRUE {

		parser.next()
		return ast.NewBooleanTrue(tok)

	} else if tok.Type == ast.T_FALSE {

		parser.next()
		return ast.NewBooleanFalse(tok)

	} else if tok.Type == ast.T_NULL {

		parser.next()
		return ast.NewNullWithToken(tok)

	} else if tok.Type == ast.T_FUNCTION_NAME {

		var fcall = parser.ParseFunctionCall()
		return ast.Expression(fcall)

	} else if tok.Type == ast.T_VARIABLE {

		return parser.ParseVariable()

	} else if tok.Type == ast.T_IDENT {

		var tok2 = parser.peekBy(2)
		if tok2 != nil && tok2.Type == ast.T_PAREN_OPEN {
			var fcall = parser.ParseFunctionCall()
			return ast.Expression(fcall)
		}

		parser.next()
		return ast.Expression(ast.NewStringWithToken(tok))

	} else if tok.Type == ast.T_HEX_COLOR {

		parser.next()
		return ast.Expression(ast.NewHexColorFromToken(tok))

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		// reduce number
		var number = parser.ParseNumber()
		return ast.Expression(number)

	} else {

		return nil
	}
	return nil
}

func (parser *Parser) ParseTerm() ast.Expression {
	debug("ParseTerm at %d", parser.Pos)
	var pos = parser.Pos
	var factor = parser.ParseFactor()
	if factor == nil {
		parser.restore(pos)
		return nil
	}

	// see if the next token is '*' or '/'
	var tok = parser.peek()
	if tok.Type == ast.T_MUL || tok.Type == ast.T_DIV {
		parser.next()
		if term := parser.ParseTerm(); term != nil {
			if tok.Type == ast.T_MUL {
				return ast.NewBinaryExpression(ast.NewOpWithToken(tok), factor, term, false)
			} else if tok.Type == ast.T_DIV {
				return ast.NewBinaryExpression(ast.NewOpWithToken(tok), factor, term, false)
			}
		} else {
			panic("Unexpected token after * and /")
		}
	}
	return factor
}

/**

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};

*/
func (parser *Parser) ParseExpression(inParenthesis bool) ast.Expression {
	var pos = parser.Pos
	debug("ParseExpression")

	// plus or minus. This creates an unary expression that holds the later term.
	// this is for:  +3 or -4
	var tok = parser.peek()
	var expr ast.Expression = nil
	if tok.Type == ast.T_PLUS || tok.Type == ast.T_MINUS {
		parser.next()
		if term := parser.ParseTerm(); term != nil {
			expr = ast.NewUnaryExpression(ast.NewOpWithToken(tok), term)

			if uexpr, ok := expr.(*ast.UnaryExpression); ok {

				// if it's evaluatable just return the evaluated value.
				if val := runtime.EvaluateUnaryExpression(uexpr, nil); val != nil {
					expr = ast.Expression(val)
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
		debug("ParseExpression failed, got %+v, restoring to %d", expr, pos)
		parser.restore(pos)
		return nil
	}

	var rightTok = parser.peek()
	for rightTok.Type == ast.T_PLUS || rightTok.Type == ast.T_MINUS || rightTok.Type == ast.T_LITERAL_CONCAT {
		// accept plus or minus
		parser.next()

		if rightTerm := parser.ParseTerm(); rightTerm != nil {
			// XXX: check parenthesis
			var bexpr = ast.NewBinaryExpression(ast.NewOpWithToken(rightTok), expr, rightTerm, inParenthesis)

			if val := runtime.EvaluateBinaryExpression(bexpr, nil); val != nil {

				expr = ast.Expression(val)

			} else {
				// wrap the existing expression with the new binary expression object
				expr = ast.Expression(bexpr)
			}
		} else {
			panic("right term is not parseable.")
		}
		rightTok = parser.peek()
	}
	return expr
}

func (parser *Parser) ParseMap() ast.Expression {
	var pos = parser.Pos
	var tok = parser.next()
	// since it's not started with '(', it's not map
	if tok.Type != ast.T_PAREN_OPEN {
		parser.restore(pos)
		return nil
	}

	tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var keyExpr = parser.ParseExpression(false)
		if keyExpr == nil {
			parser.restore(pos)
			return nil
		}

		if parser.accept(ast.T_COLON) == nil {
			parser.restore(pos)
			return nil
		}

		var valueExpr = parser.ParseExpression(false)
		if valueExpr == nil {
			parser.restore(pos)
			return nil
		}

		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			parser.next()
			tok = parser.peek()
		}
	}
	return nil
}

func (parser *Parser) ParseString() ast.Expression {
	var tok = parser.peek()

	if tok.Type == ast.T_QQ_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('"', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_Q_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('\'', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_IDENT {

		tok = parser.next()
		return ast.Expression(ast.NewStringWithToken(tok))

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	}
	return nil
}

func (parser *Parser) ParseInterp() ast.Expression {
	debug("ParseInterp at %d", parser.Pos)
	var startTok = parser.peek()

	if startTok.Type != ast.T_INTERPOLATION_START {
		return nil
	}

	parser.accept(ast.T_INTERPOLATION_START)
	var innerExpr = parser.ParseExpression(true)
	var endTok = parser.expect(ast.T_INTERPOLATION_END)
	var interp = ast.NewInterpolation(innerExpr, startTok, endTok)
	return interp
}

func (parser *Parser) ParseValueStrict() ast.Expression {
	var pos = parser.Pos

	var tok = parser.peek()
	if tok.Type == ast.T_PAREN_OPEN {
		if mapValue := parser.ParseMap(); mapValue != nil {
			return mapValue
		}
		parser.restore(pos)

		if listValue := parser.ParseList(); listValue != nil {
			return listValue
		}
		parser.restore(pos)
	}
	return parser.ParseExpression(true)
}

/**
Parse Value loosely

This parse method allows space/comma separated tokens of list.

To parse mixin argument or function argument, we only allow comma-separated list inside the parenthesis.

@param stopTokType ast.TokenType

The stop token is used from variable assignment expression,
We expect ';' semicolon at the end of expression to avoid the ambiguity of list, map and expression.
*/
func (parser *Parser) ParseValue(stopTokType ast.TokenType) ast.Expression {
	debug("ParseValue")
	var pos = parser.Pos

	// try parse map
	debug("ParseMap")
	if mapValue := parser.ParseMap(); mapValue != nil {
		var tok = parser.peek()

		if stopTokType == 0 || tok.Type == stopTokType {
			debug("OK List")
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
	debug("ParseExpression trying", pos)

	if expr := parser.ParseExpression(false); expr != nil {
		var tok = parser.peek()
		for tok.Type == ast.T_LITERAL_CONCAT {
			parser.accept(ast.T_LITERAL_CONCAT)

			var rightExpr = parser.ParseExpression(false)
			if rightExpr == nil {
				panic("Expecting expression or ident after the literal concat operator.")
			}
			expr = ast.NewLiteralConcat(expr, rightExpr)
			tok = parser.peek()
		}

		// Check if the expression is reduce-able
		// For now, division looks like CSS slash at the first level, should be string.
		if runtime.CanReduceExpression(expr) {
			if reducedExpr, ok := runtime.ReduceExpression(expr); ok {
				return reducedExpr
			}
		} else {
			return runtime.EvaluateExpression(expr, nil)
		}

		// if we can't evaluate the value, just return the expression tree
		return expr
	}
	return nil
}

func (parser *Parser) ParseList() ast.Expression {
	debug("ParseList at %d", parser.Pos)
	var pos = parser.Pos
	var list = parser.ParseCommaSepList()
	if list == nil {
		debug("ParseList failed")
		parser.restore(pos)
		return nil
	}
	return list
}

func (parser *Parser) ParseCommaSepList() ast.Expression {
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

		return list.Expressions[0]

	}
	return list
}

func (parser *Parser) ParseVariable() *ast.Variable {
	var pos = parser.Pos
	var tok = parser.next()
	if tok.Type != ast.T_VARIABLE {
		parser.restore(pos)
		return nil
	}
	return ast.NewVariableWithToken(tok)
}

func (parser *Parser) ParseVariableAssignment() ast.Statement {
	var pos = parser.Pos

	var variable = parser.ParseVariable()
	if variable == nil {
		parser.restore(pos)
		return nil
	}

	// skip ":", T_COLON token
	if parser.accept(ast.T_COLON) == nil {
		panic("Expecting colon after variable name")
	}

	// Expecting semicolon at the end of the statement
	var expr = parser.ParseValue(ast.T_SEMICOLON)
	if expr == nil {
		panic("Expecting value after variable assignment.")
	}

	if ruleset := parser.Context.TopRuleSet(); ruleset != nil {
		ruleset.Block.SymTable.Set(variable.Name, expr)
	} else if parser.Context.GlobalBlock != nil {
		parser.Context.GlobalBlock.SymTable.Set(variable.Name, expr)
	}

	var stm = ast.NewVariableAssignment(variable, expr)

	parser.ParseFlags(stm)

	parser.accept(ast.T_SEMICOLON)
	return stm
}

func (parser *Parser) ParseFlags(stm *ast.VariableAssignment) {
	var tok = parser.peek()
	for tok.IsFlagKeyword() {
		parser.next()

		switch tok.Type {
		case ast.T_DEFAULT:
			stm.Default = true
		case ast.T_OPTIONAL:
			stm.Optional = true
		case ast.T_IMPORTANT:
			stm.Important = true
		case ast.T_GLOBAL:
			stm.Global = true
		}
		tok = parser.peek()
	}
}

func (parser *Parser) ParseSpaceSepList() ast.Expression {
	debug("ParseSpaceSepList at %d", parser.Pos)

	var list = ast.NewSpaceSepList()
	list.Separator = " "

	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_OPEN {
		parser.next()
		if sublist := parser.ParseCommaSepList(); sublist != nil {
			list.Append(sublist)
		}
		parser.expect(ast.T_PAREN_CLOSE)
	}

	tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_CLOSE {
		var subexpr = parser.ParseExpression(true)
		if subexpr != nil {
			debug("Parsed Expression: %+v", subexpr)
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
		return list.Expressions[0]
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

func (parser *Parser) ParsePropertyName() ast.Expression {
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

func (parser *Parser) ParsePropertyNameToken() ast.Expression {
	var tok = parser.peek()
	if tok.Type == ast.T_PROPERTY_NAME_TOKEN {
		parser.next()
		return ast.NewIdentWithToken(tok)
	} else if tok.Type == ast.T_INTERPOLATION_START {
		return parser.ParseInterpolation()
	}
	return nil
}

func (parser *Parser) ParseInterpolation() ast.Expression {
	debug("ParseInterpolation")
	var startToken *ast.Token
	if startToken = parser.accept(ast.T_INTERPOLATION_START); startToken == nil {
		return nil
	}
	var expr = parser.ParseExpression(true)
	var endToken = parser.expect(ast.T_INTERPOLATION_END)
	return ast.NewInterpolation(expr, startToken, endToken)
}

func (parser *Parser) ParseDeclaration() ast.Statement {
	return nil
}

func (parser *Parser) ParseDeclarationBlock() *ast.DeclarationBlock {
	var declBlock = ast.DeclarationBlock{}
	var parentRuleSet = parser.Context.TopRuleSet()

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
				parser.ParseDeclarationBlock()
			}

			if parser.accept(ast.T_SEMICOLON) == nil {
				if tok3 := parser.peek(); tok3.Type == ast.T_BRACE_CLOSE {
					// normal break
					break
				} else {
					panic("missing semicolon after the property value.")
				}
			}

		} else if stm := parser.ParseStatement(); stm != nil {

		} else {
			panic(fmt.Errorf("Parse failed at token %s", tok))
		}
		tok = parser.peek()

	}
	parser.expect(ast.T_BRACE_CLOSE)
	return &declBlock
}

func (parser *Parser) ParseCharsetStatement() ast.Statement {
	parser.accept(ast.T_CHARSET)
	var tok = parser.next()
	var stm = ast.NewCharsetStatementWithToken(tok)
	parser.expect(ast.T_SEMICOLON)
	return stm
}

/*
	Media Query Syntax:
	https://developer.mozilla.org/en-US/docs/Web/Guide/CSS/Media_queries
*/
func (parser *Parser) ParseMediaQueryStatement() ast.Statement {
	// expect the '@media' token
	var stm = ast.NewMediaQueryStatement()
	parser.expect(ast.T_MEDIA)
	if list := parser.ParseMediaQueryList(); list != nil {
		stm.MediaQueryList = *list
	}
	parser.ParseBlock()
	return stm
}

func (parser *Parser) ParseMediaQueryList() *[]*ast.MediaQuery {
	var query = parser.ParseMediaQuery()
	if query == nil {
		return nil
	}

	var queries = []*ast.MediaQuery{query}

	var tok = parser.peek()
	for tok.Type == ast.T_COMMA {
		parser.next()
		if query := parser.ParseMediaQuery(); query != nil {
			queries = append(queries, query)
		}
		tok = parser.peek()
	}
	return &queries
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
		parser.next() // skip the and operator token
	}

	// parse the media expression after the media type.
	var mediaExpression = parser.ParseMediaQueryExpression()
	if mediaExpression == nil {
		if mediaType == nil {
			return nil
		}
		return ast.NewMediaQuery(mediaType, mediaExpression)
	}

	// @media query only allows AND operator here..
	var tok = parser.peek()
	for tok.Type == ast.T_LOGICAL_AND {
		parser.next()
		// parse another mediq query expression
		var expr2 = parser.ParseMediaQueryExpression()
		mediaExpression = ast.NewBinaryExpression(ast.NewOpWithToken(tok), mediaExpression, expr2, false)
		tok = parser.peek()
	}
	return ast.NewMediaQuery(mediaType, mediaExpression)
}

/*
ParseMediaType returns Ident Node or UnaryExpression as ast.Expression
*/
func (parser *Parser) ParseMediaType() ast.Expression {
	var tok = parser.peek()
	if tok.Type == ast.T_LOGICAL_NOT {
		parser.next()

		var mediaType = parser.expect(ast.T_IDENT)
		return ast.NewUnaryExpression(ast.NewOpWithToken(tok), mediaType)

	} else if tok.Type == ast.T_ONLY {
		parser.next()

		var mediaType = parser.expect(ast.T_IDENT)
		return ast.NewUnaryExpression(ast.NewOpWithToken(tok), mediaType)
	}

	// expecting media type token (it will be T_IDENT)
	tok = parser.peek()
	if tok.Type == ast.T_IDENT {
		parser.next()
		return ast.NewIdentWithToken(tok)
	}

	// parse media type fail
	return nil
}

/*
An media query expression must start with a '(' and ends with ')'
*/
func (parser *Parser) ParseMediaQueryExpression() ast.Expression {

	// it's not an media query expression
	if parser.accept(ast.T_PAREN_OPEN) == nil {
		return nil
	}

	var featureExpr = parser.ParseExpression(false)
	var feature = ast.NewMediaFeature(featureExpr, nil)

	// if the next token is a colon, then we expect a feature value
	// after the colon.
	var tok = parser.peek()
	if tok.Type == ast.T_COLON {
		parser.next()
		feature.Value = parser.ParseExpression(false)
	}
	parser.expect(ast.T_PAREN_CLOSE)
	return feature
}

func (parser *Parser) ParseWhileStatement() ast.Statement {
	parser.expect(ast.T_WHILE)
	var condition = parser.ParseCondition()
	var block = parser.ParseBlock()
	return ast.NewWhileStatement(condition, block)
}

/*
Parse the SASS @for statement.

	@for $var from <start> to <end> {  }

	@for $var from <start> through <end> {  }

@see http://sass-lang.com/documentation/file.SASS_REFERENCE.html#_10
*/
func (parser *Parser) ParseForStatement() ast.Statement {
	parser.expect(ast.T_FOR)

	// get the variable token
	var variable = parser.ParseVariable()
	var stm = ast.NewForStatement(variable)

	if parser.accept(ast.T_FOR_FROM) != nil {

		var fromExpr = parser.ParseExpression(true)
		if reducedExpr, ok := runtime.ReduceExpression(fromExpr); ok {
			fromExpr = reducedExpr
		}
		stm.From = fromExpr

		// "through" or "to"
		var tok = parser.next()

		if tok.Type != ast.T_FOR_THROUGH && tok.Type != ast.T_FOR_TO {
			panic("Expecting 'through' or 'to' of range syntax.")
		}

		var endExpr = parser.ParseExpression(true)
		if reducedExpr, ok := runtime.ReduceExpression(endExpr); ok {
			endExpr = reducedExpr
		}

		if tok.Type == ast.T_FOR_THROUGH {

			stm.Through = endExpr

		} else if tok.Type == ast.T_FOR_TO {

			stm.To = endExpr

		}

	} else if parser.accept(ast.T_FOR_IN) != nil {

		var fromExpr = parser.ParseExpression(true)
		if reducedExpr, ok := runtime.ReduceExpression(fromExpr); ok {
			fromExpr = reducedExpr
		}
		stm.From = fromExpr

		parser.expect(ast.T_RANGE)

		var endExpr = parser.ParseExpression(true)
		if reducedExpr, ok := runtime.ReduceExpression(endExpr); ok {
			endExpr = reducedExpr
		}

		stm.To = endExpr
	}

	if b := parser.ParseBlock(); b != nil {
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
func (parser *Parser) ParseImportStatement() ast.Statement {
	// skip the ast.T_IMPORT token
	parser.expect(ast.T_IMPORT)

	// Create the import statement node
	var stm = ast.NewImportStatement()

	var tok = parser.peek()
	// expecting url(..)
	if tok.Type == ast.T_IDENT {
		parser.advance()

		if tok.Str != "url" {
			panic("invalid function for @import statement.")
		}

		if tok = parser.next(); tok.Type != ast.T_PAREN_OPEN {
			panic("expecting parenthesis after url")
		}

		tok = parser.next()
		stm.Url = ast.Url(tok.Str)

		if tok = parser.next(); tok.Type != ast.T_PAREN_CLOSE {
			panic("expecting parenthesis after url")
		}

	} else if tok.IsString() {
		parser.advance()
		stm.Url = ast.RelativeUrl(tok.Str)
	}

	parser.ParseMediaQueryList()

	// must be ast.T_SEMICOLON
	tok = parser.next()
	if tok.Type != ast.T_SEMICOLON {
		panic(ParserError{";", tok.Str})
	}
	return stm
}

func (parser *Parser) ParseReturnStatement() ast.Statement {
	var returnTok = parser.expect(ast.T_RETURN)
	var valueExpr = parser.ParseExpression(true)
	parser.expect(ast.T_SEMICOLON)
	return ast.NewReturnStatementWithToken(returnTok, valueExpr)
}

func (parser *Parser) ParseFunction() ast.Statement {
	parser.expect(ast.T_FUNCTION)
	var identTok = parser.expect(ast.T_FUNCTION_NAME)
	var args = parser.ParseFunctionPrototype()
	var fun = ast.NewFunctionWithToken(identTok)
	fun.ArgumentList = args
	fun.Block = parser.ParseBlock()
	return fun
}

func (parser *Parser) ParseMixinStatement() ast.Statement {
	var mixinTok = parser.expect(ast.T_MIXIN)
	var stm = ast.NewMixinStatementWithToken(mixinTok)

	var tok = parser.next()

	// Mixin without parameters
	if tok.Type == ast.T_IDENT {

		stm.Ident = tok

	} else if tok.Type == ast.T_FUNCTION_NAME {

		stm.Ident = tok

		stm.ArgumentList = parser.ParseFunctionPrototype()
	}

	stm.Block = parser.ParseDeclarationBlock()
	return stm
}

func (parser *Parser) ParseFunctionPrototypeArgument() *ast.Argument {
	debug("ParseFunctionPrototypeArgument")

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

func (parser *Parser) ParseFunctionPrototype() *ast.ArgumentList {
	debug("ParseFunctionPrototype")

	var args = ast.NewArgumentList()

	parser.expect(ast.T_PAREN_OPEN)
	var tok = parser.peek()
	for tok.Type != ast.T_PAREN_CLOSE {
		var arg *ast.Argument = nil
		if arg = parser.ParseFunctionPrototypeArgument(); arg != nil {
			args.Append(arg)
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

func (parser *Parser) ParseIncludeStatement() ast.Statement {
	var tok = parser.expect(ast.T_INCLUDE)
	var stm = ast.NewIncludeStatementWithToken(tok)

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
		stm.ContentBlock = parser.ParseDeclarationBlock()
	}

	parser.expect(ast.T_SEMICOLON)
	return stm
}

/*
@content directive is only allowed in mixin block
*/
func (parser *Parser) ParseContentStatement() ast.Statement {
	var tok = parser.expect(ast.T_CONTENT)
	parser.expect(ast.T_SEMICOLON)
	return ast.NewContentStatementWithToken(tok)
}
