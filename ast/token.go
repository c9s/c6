package ast

import "fmt"

//go:generate stringer -type=TokenType token.go
type TokenType int

const LF = '\r'
const CR = '\n'

type KeywordTokenMap map[string]TokenType

type KeywordToken struct {
	Keyword   string
	TokenType TokenType
}

type KeywordTokenList []KeywordToken

var FlagTokenList = KeywordTokenList{
	KeywordToken{"!global", T_FLAG_GLOBAL},
	KeywordToken{"!default", T_FLAG_DEFAULT},
	KeywordToken{"!important", T_FLAG_IMPORTANT},
	KeywordToken{"!optional", T_FLAG_OPTIONAL},
	KeywordToken{"!constant", T_FLAG_CONSTANT},
}

var ForRangeKeywordTokenList = KeywordTokenList{
	KeywordToken{"from", T_FOR_FROM},
	KeywordToken{"through", T_FOR_THROUGH},
	KeywordToken{"to", T_FOR_TO},
	KeywordToken{"in", T_FOR_IN},
}

// TODO: sort by frequency
var KeywordList = []KeywordToken{
	KeywordToken{"@else if", T_ELSE_IF},
	KeywordToken{"@else", T_ELSE},
	KeywordToken{"@if", T_IF},
	KeywordToken{"@import", T_IMPORT},
	KeywordToken{"@media", T_MEDIA},
	KeywordToken{"@page", T_PAGE},
	KeywordToken{"@return", T_RETURN},
	KeywordToken{"@each", T_EACH},
	KeywordToken{"@when", T_WHEN},
	KeywordToken{"@include", T_INCLUDE},
	KeywordToken{"@function", T_FUNCTION},
	KeywordToken{"@mixin", T_MIXIN},
	KeywordToken{"@font-face", T_FONT_FACE},
	KeywordToken{"@for", T_FOR},
	KeywordToken{"@error", T_ERROR},
	KeywordToken{"@warn", T_WARN},
	KeywordToken{"@info", T_INFO},
	KeywordToken{"@while", T_WHILE},
	KeywordToken{"@content", T_CONTENT},
	KeywordToken{"@extend", T_EXTEND},
	KeywordToken{"@namespace", T_NAMESPACE},
}

var ExprTokenList = KeywordTokenList{
	KeywordToken{"true", T_TRUE},
	KeywordToken{"false", T_FALSE},
	KeywordToken{"null", T_NULL},
	KeywordToken{"only", T_ONLY}, // used in media query
	KeywordToken{"and", T_LOGICAL_AND},
	KeywordToken{"not", T_LOGICAL_NOT},
	KeywordToken{"or", T_LOGICAL_OR},
	KeywordToken{"xor", T_LOGICAL_XOR},
	KeywordToken{"odd", T_ODD},
	KeywordToken{"even", T_EVEN},
}

var UnitTokenList = KeywordTokenList{
	KeywordToken{"px", T_UNIT_PX},
	KeywordToken{"pt", T_UNIT_PT},
	KeywordToken{"pc", T_UNIT_PC},
	KeywordToken{"em", T_UNIT_EM},
	KeywordToken{"cm", T_UNIT_CM},
	KeywordToken{"ex", T_UNIT_EX},
	KeywordToken{"ch", T_UNIT_CH},
	KeywordToken{"in", T_UNIT_IN},
	KeywordToken{"mm", T_UNIT_MM},
	KeywordToken{"rem", T_UNIT_REM},
	KeywordToken{"vh", T_UNIT_VH},
	KeywordToken{"vw", T_UNIT_VW},

	KeywordToken{"Hz", T_UNIT_HZ},
	KeywordToken{"kHz", T_UNIT_KHZ},

	KeywordToken{"vmin", T_UNIT_VMIN},
	KeywordToken{"vmax", T_UNIT_VMAX},
	KeywordToken{"deg", T_UNIT_DEG},
	KeywordToken{"grad", T_UNIT_GRAD},
	KeywordToken{"rad", T_UNIT_RAD},
	KeywordToken{"turn", T_UNIT_TURN},
	KeywordToken{"dpi", T_UNIT_DPI},
	KeywordToken{"dpcm", T_UNIT_DPCM},
	KeywordToken{"dppx", T_UNIT_DPPX},
	KeywordToken{"s", T_UNIT_SECOND},
	KeywordToken{"ms", T_UNIT_MILLISECOND},
	KeywordToken{"%", T_UNIT_PERCENT},
}

type Token struct {
	Type                  TokenType
	Str                   string
	Pos                   int
	Line                  int
	LineOffset            int
	ContainsInterpolation bool
}

type TokenStream chan *Token

/**
Implement the stringer interface
*/
func (tok Token) String() string {
	return fmt.Sprintf("'%s' (%s) at line %d, offset %d", tok.Str, tok.Type, tok.Line, tok.Pos)
}

func (tok Token) IsString() bool {
	return tok.Type == T_QQ_STRING || tok.Type == T_Q_STRING || tok.Type == T_UNQUOTE_STRING
}

func (tok Token) IsSelectorCombinator() bool {
	return tok.Type == T_ADJACENT_SIBLING_COMBINATOR ||
		tok.Type == T_CHILD_COMBINATOR ||
		tok.Type == T_GENERAL_SIBLING_COMBINATOR ||
		tok.Type == T_DESCENDANT_COMBINATOR ||
		tok.Type == T_COMMA
}

func (tok Token) IsSelector() bool {
	switch tok.Type {
	case T_TYPE_SELECTOR, T_UNIVERSAL_SELECTOR, T_ID_SELECTOR,
		T_CLASS_SELECTOR, T_PARENT_SELECTOR,
		T_ADJACENT_SIBLING_COMBINATOR, T_GENERAL_SIBLING_COMBINATOR,
		T_CHILD_COMBINATOR, T_DESCENDANT_COMBINATOR,
		T_PSEUDO_SELECTOR,
		T_FUNCTIONAL_PSEUDO,
		T_BRACKET_OPEN: // '[' is the first token of attribute selector.
		return true
	}
	return false
}

func (tok Token) IsAttributeMatchOperator() bool {
	switch tok.Type {
	case T_ATTR_EQUAL,
		T_INCLUDE_MATCH,   // for '~='
		T_PREFIX_MATCH,    // for '^='
		T_DASH_MATCH,      // for '|='
		T_SUFFIX_MATCH,    // for '$='
		T_SUBSTRING_MATCH: // for '*='
		return true
	}
	return false
}

func (tok Token) IsFlagKeyword() bool {
	switch tok.Type {
	case T_FLAG_DEFAULT, T_FLAG_OPTIONAL, T_FLAG_GLOBAL, T_FLAG_IMPORTANT:
		return true
	}
	return false
}

func (tok Token) IsUnit() bool {
	switch tok.Type {
	case T_UNIT_NONE, T_UNIT_PERCENT, T_UNIT_SECOND, T_UNIT_MILLISECOND,
		T_UNIT_EM, T_UNIT_EX, T_UNIT_CH, T_UNIT_REM, T_UNIT_CM, T_UNIT_IN,
		T_UNIT_MM, T_UNIT_PC, T_UNIT_PT, T_UNIT_PX, T_UNIT_VH, T_UNIT_VW,
		T_UNIT_VMIN, T_UNIT_VMAX, T_UNIT_HZ, T_UNIT_KHZ, T_UNIT_DPI, T_UNIT_DPCM,
		T_UNIT_DPPX, T_UNIT_DEG, T_UNIT_GRAD, T_UNIT_RAD, T_UNIT_TURN, T_UNIT_OTHERS:
		return true
	}
	return false
}

func (tok Token) IsComparisonOperator() bool {
	switch tok.Type {
	case T_EQUAL, T_UNEQUAL, T_GT, T_LT, T_GE, T_LE:
		return true
	}
	return false
}

func (tok Token) IsOneOfTypes(types []TokenType) bool {
	for _, t := range types {
		if tok.Type == t {
			return true
		}
	}
	return false
}

const (
	T_SPACE TokenType = iota
	T_COMMENT_LINE
	T_COMMENT_BLOCK
	T_SEMICOLON
	T_COMMA
	T_IDENT
	T_URL
	T_MEDIA
	T_PAGE

	T_TRUE
	T_FALSE
	T_NULL
	T_ONLY
	T_ODD  // for nth-child(odd)
	T_EVEN // for nth-child(even)
	T_N    // for nth-child(3n+3)

	T_MS_PARAM_NAME
	T_FUNCTION_NAME

	// selector tokens
	T_ID_SELECTOR
	T_CLASS_SELECTOR
	T_TYPE_SELECTOR
	T_UNIVERSAL_SELECTOR
	T_PARENT_SELECTOR   // SASS parent selector
	T_PSEUDO_SELECTOR   // :hover, :visited , ...
	T_FUNCTIONAL_PSEUDO // lang(...), nth(...)

	/*
		An interpolation selector token presents one or two more selector strings,
		which may contains an expression that change the type of the selector.
	*/
	T_INTERPOLATION_SELECTOR // selector with interpolation: '#{ ... }'

	/*
		The literal concat means we would concat two string without quotes.
		This is used for concating strings or expression with interpolation sections.
	*/
	T_LITERAL_CONCAT // used to concat selectors and interpolation

	/*
		This is for normal string concat
	*/
	T_CONCAT

	// for Microsoft 'progid:' token, we don't have choice.
	T_MS_PROGID

	// Selector relationship
	T_DESCENDANT_COMBINATOR       // E ' ' F
	T_CHILD_COMBINATOR            // E '>' F
	T_ADJACENT_SIBLING_COMBINATOR // E '+' F
	T_GENERAL_SIBLING_COMBINATOR  // E '~' F

	T_UNICODE_RANGE

	T_IF       // @if
	T_ELSE     // @else
	T_ELSE_IF  // @else if
	T_INCLUDE  // for @include
	T_EACH     // for @each
	T_WHEN     // for @when
	T_MIXIN    // for @mixin
	T_EXTEND   // for @extend
	T_FUNCTION // for @function
	T_WARN     // @warn
	T_ERROR    // @error
	T_INFO     // @info
	T_FOR
	T_FOR_FROM
	T_FOR_THROUGH
	T_FOR_TO
	T_FOR_IN  // for @for ... in
	T_WHILE   // for @while
	T_RETURN  // for @return
	T_RANGE   // for '..'
	T_CONTENT // for '@content'

	// Flag token types
	T_FLAG_GLOBAL
	T_FLAG_DEFAULT
	T_FLAG_IMPORTANT
	T_FLAG_OPTIONAL

	T_FONT_FACE
	T_NAMESPACE

	T_LOGICAL_NOT // 'not' used in conditions
	T_LOGICAL_OR  // 'or' used in conditions query
	T_LOGICAL_AND // 'and' used in conditions query
	T_LOGICAL_XOR

	/*
		expression operators
	*/
	T_NOP
	T_PLUS  // for '+'
	T_DIV   // for '-'
	T_MUL   // for '*'
	T_MINUS // for '-'
	T_MOD   // for '%'
	T_BRACE_OPEN
	T_BRACE_CLOSE
	T_LANG_CODE // 'en', 'fr', 'fr-ca'
	T_BRACKET_OPEN
	T_ATTRIBUTE_NAME
	T_BRACKET_CLOSE

	T_EQUAL   // for '=='
	T_UNEQUAL // for '!='
	T_GT      // greater than for '>'
	T_LT      // less than for '<'
	T_GE      // for '>='
	T_LE      // for '<='

	T_ASSIGN     // for '='
	T_ATTR_EQUAL // for '=' inside attribute selecctors

	T_INCLUDE_MATCH   // for '~='
	T_PREFIX_MATCH    // for '^='
	T_DASH_MATCH      // for '|='
	T_SUFFIX_MATCH    // for '$='
	T_SUBSTRING_MATCH // for '*='

	T_VARIABLE                  // for $[a-z]*
	T_VARIABLE_LENGTH_ARGUMENTS // for '...'

	T_IMPORT
	T_AT_RULE

	T_CHARSET
	T_QQ_STRING
	T_Q_STRING
	T_UNQUOTE_STRING
	T_PAREN_OPEN
	T_PAREN_CLOSE
	T_FLAG_CONSTANT
	T_INTEGER
	T_FLOAT

	T_CDO // for <!--
	T_CDC // for -->

	/*
		unit tokens
	*/
	T_UNIT_NONE
	T_UNIT_OTHERS
	T_UNIT_PERCENT

	/*
		Time Unit
		@see https://developer.mozilla.org/zh-TW/docs/Web/CSS/time
	*/
	T_UNIT_SECOND
	T_UNIT_MILLISECOND

	/*
		Length
	*/

	/*
		Length Unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/length
	*/
	T_UNIT_EM
	T_UNIT_EX
	T_UNIT_CH
	T_UNIT_REM

	// Absolute length
	T_UNIT_CM
	T_UNIT_IN
	T_UNIT_MM
	T_UNIT_PC
	T_UNIT_PT
	T_UNIT_PX

	// Viewport-percentage lengths
	T_UNIT_VH
	T_UNIT_VW
	T_UNIT_VMIN
	T_UNIT_VMAX

	/*
		Frequency unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/frequency
	*/
	T_UNIT_HZ
	T_UNIT_KHZ

	/*
		Resolution Unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/resolution
	*/
	T_UNIT_DPI
	T_UNIT_DPCM
	T_UNIT_DPPX

	/*
		Angle
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/angle
	*/
	T_UNIT_DEG
	T_UNIT_GRAD
	T_UNIT_RAD
	T_UNIT_TURN

	T_PROPERTY_NAME_TOKEN
	T_PROPERTY_VALUE
	T_HEX_COLOR
	T_COLON
	T_INTERPOLATION_START
	T_INTERPOLATION_INNER
	T_INTERPOLATION_END
)
