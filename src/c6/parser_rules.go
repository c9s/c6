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
	parser.expect(ast.T_BRACE_START)
	var block = ast.NewBlock()
	block.Statements = parser.ParseStatements()
	parser.expect(ast.T_BRACE_END)
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

	if token.Type == ast.T_IMPORT {

		return parser.ParseImportStatement()

	} else if token.Type == ast.T_CHARSET {

		return parser.ParseCharsetStatement()

	} else if token.Type == ast.T_MEDIA {

		return parser.ParseMediaQueryStatement()

	} else if token.Type == ast.T_VARIABLE {

		return parser.ParseVariableAssignment()

	} else if token.Type == ast.T_IF {

		return parser.ParseIfStatement()

	} else if token.IsSelector() {

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
	if tok.Type == ast.T_PAREN_START {
		parser.accept(ast.T_PAREN_START)
		expr = parser.ParseLogicExpression()
		parser.expect(ast.T_PAREN_END)
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

func (parser *Parser) ParseRuleSet() ast.Statement {
	var tok = parser.next()
	var parentRuleSet = parser.Context.TopRuleSet()

	var ruleset = ast.NewRuleSet()
	parser.Context.PushRuleSet(ruleset)

	for tok.IsSelector() {

		switch tok.Type {

		case ast.T_TYPE_SELECTOR:
			sel := &ast.TypeSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_UNIVERSAL_SELECTOR:
			sel := &ast.UniversalSelector{}
			ruleset.AppendSelector(sel)

		case ast.T_ID_SELECTOR:
			sel := &ast.IdSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_CLASS_SELECTOR:
			sel := &ast.ClassSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_PARENT_SELECTOR:

			sel := &ast.ParentSelector{parentRuleSet}
			ruleset.AppendSelector(sel)

		case ast.T_PSEUDO_SELECTOR:

			sel := &ast.PseudoSelector{tok.Str, ""}
			if nextTok := parser.peek(); nextTok.Type == ast.T_LANG_CODE {
				sel.C = nextTok.Str
			}
			ruleset.AppendSelector(sel)

		case ast.T_ADJACENT_SIBLING_COMBINATOR:

			ruleset.AppendSelector(&ast.AdjacentSelector{})

		case ast.T_CHILD_COMBINATOR:

			ruleset.AppendSelector(&ast.ChildSelector{})

		case ast.T_DESCENDANT_COMBINATOR:

			ruleset.AppendSelector(&ast.DescendantSelector{})

		default:
			panic(fmt.Errorf("Unexpected selector token: %+v", tok))
		}
		tok = parser.next()
	}
	parser.backup()

	// parse declaration block
	ruleset.Block = parser.ParseDeclarationBlock()

	// pop the ruleset from stack
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

	var fcall = ast.NewFunctionCall(identTok)

	parser.expect(ast.T_PAREN_START)

	var argTok = parser.peek()
	for argTok.Type != ast.T_PAREN_END {
		var arg = parser.ParseFactor()
		fcall.AppendArgument(arg)
		debug("ParseFunctionCall => arg: %+v", arg)

		argTok = parser.peek()
		if argTok.Type == ast.T_COMMA {
			parser.next() // skip comma
			argTok = parser.peek()
		} else if argTok.Type == ast.T_PAREN_END {
			parser.next() // consume ')'
			break
		}
	}
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

	if tok.Type == ast.T_PAREN_START {
		parser.expect(ast.T_PAREN_START)
		var expr = parser.ParseExpression(true)
		parser.expect(ast.T_PAREN_END)
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

	} else if tok.Type == ast.T_IDENT {

		parser.next()
		return ast.Expression(ast.NewStringWithToken(tok))

	} else if tok.Type == ast.T_HEX_COLOR {

		parser.next()
		return ast.Expression(ast.NewHexColorFromToken(tok))

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		// reduce number
		var number = parser.ParseNumber()
		return ast.Expression(number)

	} else if tok.Type == ast.T_FUNCTION_NAME {

		var fcall = parser.ParseFunctionCall()
		return ast.Expression(fcall)

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
	if tok.Type != ast.T_PAREN_START {
		parser.restore(pos)
		return nil
	}

	tok = parser.peek()
	for tok.Type != ast.T_PAREN_END {
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

/**
The stop token is used from variable assignment expression,
 we expect ';' semicolon at the end of expression to avoid the ambiguity of list, map and expression.
*/
func (parser *Parser) ParseValue(stopTokType ast.TokenType) ast.Expression {
	debug("ParseValue")
	var pos = parser.Pos

	// try parse map
	debug("Trying Map")
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
	for tok.Type != ast.T_COMMA && tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {

		// when the syntax start with a '(', it could be a list or map.
		if tok.Type == ast.T_PAREN_START {

			parser.next()
			if sublist := parser.ParseCommaSepList(); sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			}
			// allow empty list here
			parser.expect(ast.T_PAREN_END)

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
	return ast.NewVariable(tok)
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

	if tok.Type == ast.T_PAREN_START {
		parser.next()
		if sublist := parser.ParseCommaSepList(); sublist != nil {
			list.Append(sublist)
		}
		parser.expect(ast.T_PAREN_END)
	}

	tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
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
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
		var sublist = parser.ParseList()
		if sublist != nil {
			list.Append(sublist)
			debug("ParsePropertyValue list: %+v", list)
		} else {
			break
		}
		tok = parser.peek()
	}

	tok = parser.peek()
	if tok.Type == ast.T_SEMICOLON || tok.Type == ast.T_BRACE_END {
		parser.next()
	} else {
		panic(fmt.Errorf("Unexpected end of property value. Got %s", tok))
	}
	return list
}

func (parser *Parser) ParseDeclarationBlock() *ast.DeclarationBlock {
	var declBlock = ast.DeclarationBlock{}
	var parentRuleSet = parser.Context.TopRuleSet()

	parser.expect(ast.T_BRACE_START)

	var tok = parser.peek()
	for tok != nil && tok.Type != ast.T_BRACE_END {
		if tok.Type == ast.T_PROPERTY_NAME_TOKEN {
			// TODO support property interpolation here
			parser.expect(ast.T_PROPERTY_NAME_TOKEN)
			parser.expect(ast.T_COLON)

			var property = ast.NewProperty(tok)
			var valueList = parser.ParsePropertyValue(parentRuleSet, property)
			_ = valueList
			// property.Values = valueList
			declBlock.Append(property)
			_ = property

		} else {

			if stm := parser.ParseStatements(); stm == nil {
				// TODO: throw syntax error report here
				panic("Can't parse statements")
			}
		}

		tok = parser.peek()
	}

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
	} else {
		// XXX: report detail here
		panic("@media query syntax error")
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
	if parser.accept(ast.T_PAREN_START) == nil {
		return nil
	}

	var featureExpr = parser.ParseExpression(false)
	_ = featureExpr

	// if the next token is a colon, then we expect a feature value
	// after the colon.
	var tok = parser.peek()
	if tok.Type == ast.T_COLON {
		parser.next()
		var featureValue = parser.ParseExpression(false)
		_ = featureValue
	}

	parser.expect(ast.T_PAREN_END)
	return nil
}

func (parser *Parser) ParseImportStatement() ast.Statement {
	// skip the ast.T_IMPORT token
	var tok = parser.next()

	// Create the import statement node
	var stm = ast.NewImportStatement()

	tok = parser.peek()
	// expecting url(..)
	if tok.Type == ast.T_IDENT {
		parser.advance()

		if tok.Str != "url" {
			panic("invalid function for @import statement.")
		}

		if tok = parser.next(); tok.Type != ast.T_PAREN_START {
			panic("expecting parenthesis after url")
		}

		tok = parser.next()
		stm.Url = ast.Url(tok.Str)

		if tok = parser.next(); tok.Type != ast.T_PAREN_END {
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
