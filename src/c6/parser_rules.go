package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "fmt"
import "strconv"
import "c6/ast"

func (parser *Parser) ParseStatement(parentRuleSet *ast.RuleSet) ast.Statement {
	var token = parser.peek()

	if token.Type == ast.T_IMPORT {
		return parser.ParseImportStatement()
	} else if token.Type == ast.T_VARIABLE {
		return parser.ParseVariableAssignment()
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

/**
This method returns objects with ast.Number interface

works for:

	'10'
	'10' 'px'
	'10' 'em'
	'0.2' 'em'
*/
func (parser *Parser) ParseNumber() *ast.Number {
	// the number token
	var tok = parser.next()

	debug("ParseNumber => next: %s", tok)

	var tok2 = parser.peek()
	var number *ast.Number
	if tok.Type == ast.T_INTEGER {
		i, err := strconv.ParseInt(tok.Str, 10, 64)
		if err != nil {
			panic(err)
		}
		number = ast.NewNumberInt64(i, tok)

	} else {

		f, err := strconv.ParseFloat(tok.Str, 64)
		if err != nil {
			panic(err)
		}
		number = ast.NewNumber(f, tok)
	}

	if tok2.IsOneOfTypes([]ast.TokenType{ast.T_UNIT_PX, ast.T_UNIT_PT, ast.T_UNIT_CM, ast.T_UNIT_EM, ast.T_UNIT_MM, ast.T_UNIT_REM, ast.T_UNIT_DEG, ast.T_UNIT_PERCENT}) {
		// consume the unit token
		parser.next()
		number.SetUnit(ast.ConvertTokenTypeToUnitType(tok2.Type))
	}
	return number
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
	return ast.NewIdent(tok.Str, *tok)
}

/**
The ParseFactor must return an Expression interface compatible object
*/
func (parser *Parser) ParseFactor() ast.Expression {
	var tok = parser.peek()
	debug("ParseFactor => peek: %s", tok)

	if tok.Type == ast.T_PAREN_START {

		parser.expect(ast.T_PAREN_START)
		var expr = parser.ParseExpression()
		parser.expect(ast.T_PAREN_END)
		return expr

	} else if tok.Type == ast.T_INTERPOLATION_START {

		parser.expect(ast.T_INTERPOLATION_START)
		parser.ParseExpression()
		parser.expect(ast.T_INTERPOLATION_END)
		// TODO:

	} else if tok.Type == ast.T_QQ_STRING || tok.Type == ast.T_Q_STRING {

		tok = parser.next()
		var str = ast.NewString(tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		// reduce number
		var number = parser.ParseNumber()
		return ast.Expression(number)

	} else if tok.Type == ast.T_FUNCTION_NAME {

		var fcall = parser.ParseFunctionCall()
		return ast.Expression(*fcall)

	} else if tok.Type == ast.T_IDENT {

		var ident = parser.ParseIdent()
		return ast.Expression(ident)

	} else if tok.Type == ast.T_HEX_COLOR {

		panic("hex color is not implemented yet")

		// TODO: Add more incorrect cases here

	} else {

		return nil
	}
	return nil
}

func (parser *Parser) ParseTerm() ast.Expression {
	var factor = parser.ParseFactor()
	if factor == nil {
		return nil
	}

	// see if the next token is '*' or '/'
	var tok = parser.peek()
	if tok.Type == ast.T_MUL || tok.Type == ast.T_DIV {
		parser.next()
		var term = parser.ParseTerm()

		if tok.Type == ast.T_MUL {
			return ast.NewBinaryExpression(ast.OpMul, factor, term)
		} else if tok.Type == ast.T_DIV {
			return ast.NewBinaryExpression(ast.OpDiv, factor, term)
		}
	}
	return factor
}

/**

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};

*/
func (parser *Parser) ParseExpression() ast.Expression {
	if tok := parser.accept(ast.T_INTERPOLATION_START); tok != nil {
		debug("ParseExpression => accept: T_INTERPOLATION_START")

		debug("ParseExpression => ParseExpression")
		var expr = parser.ParseExpression()

		endToken := parser.expect(ast.T_INTERPOLATION_END)
		debug("ParseExpression => expect: T_INTERPOLATION_START")

		var interp = ast.NewInterpolation(expr, tok, endToken)
		return interp
	}

	// plus or minus. this creates an unary expression that holds the later term.
	// this is for:  +3 or -4
	var tok = parser.peek()
	var expr ast.Expression
	if tok.Type == ast.T_PLUS || tok.Type == ast.T_MINUS {
		parser.next()
		var term = parser.ParseTerm()
		expr = ast.NewUnaryExpression(ast.ConvertTokenTypeToOpType(tok.Type), term)
	} else {
		expr = parser.ParseTerm()
	}

	var rightTok = parser.peek()
	for rightTok.Type == ast.T_PLUS || rightTok.Type == ast.T_MINUS {
		parser.next()
		var rightTerm = parser.ParseTerm()
		expr = ast.NewBinaryExpression(ast.ConvertTokenTypeToOpType(rightTok.Type), expr, rightTerm)
		rightTok = parser.peek()
	}
	return expr
}

/*
func (parser *Parser) ParseMap() *ast.Map {
	var tok = parser.next()
	if tok.Type != ast.T_PAREN_START {
		panic("Map Syntax error: expecting '('")
	}

	tok = parser.next()
	if tok.Type != ast.T_IDENT {
		panic("Map Syntax error: expecting ident or expression after '('")
	}
	return nil
}
*/

func (parser *Parser) ParseMap() ast.Node {
	var pos = parser.Pos
	var tok = parser.next()
	// since it's not started with '(', it's not map
	if tok.Type != ast.T_PAREN_START {
		parser.restore(pos)
		return nil
	}

	tok = parser.peek()
	for tok.Type != ast.T_PAREN_END {
		var keyExpr = parser.ParseExpression()
		if keyExpr == nil {
			parser.restore(pos)
			return nil
		}

		if parser.expect(ast.T_COLON) == nil {
			parser.restore(pos)
			return nil
		}

		var valueExpr = parser.ParseExpression()
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

func (parser *Parser) ParseValue() ast.Node {
	debug("ParseValue")
	var pos = parser.Pos

	if mapValue := parser.ParseMap(); mapValue != nil {
		return mapValue
	}

	parser.restore(pos)

	var tok = parser.peek()

	// list or map starts with '('
	if tok.Type == ast.T_PAREN_START {
		var expr = parser.ParseExpression()
		_ = expr
	}

	tok = parser.peek()
	if tok.Type == ast.T_COLON {
		// it's a map
	}

	tok = parser.peek()
	if tok.Type == ast.T_PAREN_START {
		// parser.ParseMapOrList()
	} else {
		parser.ParseList()
	}
	return nil
}

func (parser *Parser) ParseList() *ast.List {
	var list = ast.NewList()
	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
		if sublist := parser.ParseCommaSepList(); sublist != nil {
			list.Append(sublist)
		}

		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			parser.next()
		}
		tok = parser.peek()
	}
	fmt.Printf("ParseList <= %s\n", list.String())
	return list
}

func (parser *Parser) ParseCommaSepList() *ast.List {
	var list = ast.NewList()
	list.Separator = ", "

	var tok = parser.peek()
	for tok.Type != ast.T_COMMA && tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {

		// when the syntax start with a '(', it could be a list or map.
		if tok.Type == ast.T_PAREN_START {
			parser.expect(ast.T_PAREN_START)

			var sublist = parser.ParseCommaSepList()

			parser.expect(ast.T_PAREN_END)

			if sublist != nil {
				list.Append(sublist)
			}
		} else {
			var sublist = parser.ParseSpaceSepList()
			if sublist != nil {
				list.Append(sublist)
			}
		}

		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			parser.next()
			tok = parser.peek()
		}
	}
	debug("Comma-separated list: %s\n", list)

	// XXX: if there is only one item in the list, we can reduce it to element.

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
		parser.restore(pos)
		return nil
	}

	var expr = parser.ParseExpression()
	if expr == nil {
		parser.restore(pos)
		return nil
	}

	parser.expect(ast.T_SEMICOLON)

	// Reduce list or map here
	return ast.NewVariableAssignment(variable, expr)
}

func (parser *Parser) ParseSpaceSepList() *ast.List {
	var list = ast.NewList()
	list.Separator = " "

	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_START {
		parser.next()
		var sublist = parser.ParseCommaSepList()
		parser.expect(ast.T_PAREN_END)

		list.Append(sublist)
	}

	tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
		var expr = parser.ParseExpression()
		if expr != nil {
			list.Append(expr)
		}
		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			break
		}
	}
	return list
}

/*
func (parser *Parser) ReduceCommaSepList() *ast.List {
	var tok = parser.peek()

	// the start of the list token
	if tok.Type == ast.T_PAREN_START {
		parser.next()
		// reduce the space separated list inside
		var grouplist = parser.ReduceList()
		_ = grouplist
	} else {
		var list = ast.NewList()
		var expr = parser.ParseExpression()
		var tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			list.Separator = ","
		} else if tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
			list.Separator = " "
		}

		// try to parse the second expression
		var expr2 = parser.ParseExpression()
		for expr2 != nil {
			var tok = parser.peek()

			// the processed token is comma, but what we've got is not a comma
			// this is for "," then " " case
			if list.Separator == "," && tok.Type != ast.T_COMMA {

			} else if list.Separator == " " && tok.Type == ast.T_COMMA {

			}

			expr2 = parser.ParseExpression()
		}
	}
	return nil
}
*/

/**
We treat the property value section as a list value, which is separated by ',' or ' '
*/
func (parser *Parser) ParsePropertyValue(parentRuleSet *ast.RuleSet, property *ast.Property) *ast.List {
	// var tok = parser.peek()
	var list = parser.ParseList()

	var tok = parser.peek()

	if tok.Type == ast.T_SEMICOLON || tok.Type == ast.T_BRACE_END {
		parser.next()
	} else {
		panic(fmt.Errorf("Unexpected end of property value. Got %s", tok))
	}
	return list
}

func (parser *Parser) ParseDeclarationBlock(parentRuleSet *ast.RuleSet) *ast.DeclarationBlock {
	var declBlock = ast.DeclarationBlock{}

	var tok = parser.next() // should be '{'
	if tok.Type != ast.T_BRACE_START {
		panic(ParserError{"{", tok.Str})
	}

	tok = parser.next()
	for tok != nil && tok.Type != ast.T_BRACE_END {

		if tok.Type == ast.T_PROPERTY_NAME_TOKEN {
			parser.expect(ast.T_COLON)

			var property = ast.NewProperty(tok)
			var valueList = parser.ParsePropertyValue(parentRuleSet, property)
			_ = valueList
			// property.Values = valueList
			declBlock.Append(property)
			_ = property

		} else if tok.IsSelector() {
			// parse subrule
			panic("subselector unimplemented")
		} else {
			panic("unexpected token")
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
