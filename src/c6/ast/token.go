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

var FlagTokenMap = KeywordTokenMap{
	"!global":    T_GLOBAL,
	"!default":   T_DEFAULT,
	"!important": T_IMPORTANT,
	"!optional":  T_OPTIONAL,
}

var ForRangeKeywordTokenMap = KeywordTokenMap{
	"from":    T_FOR_FROM,
	"through": T_FOR_THROUGH,
	"to":      T_FOR_TO,
	"in":      T_FOR_IN,
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
	KeywordToken{"@while", T_WHILE},
	KeywordToken{"@content", T_CONTENT},
}

var ExprTokenMap = KeywordTokenMap{
	"true":  T_TRUE,
	"false": T_FALSE,
	"null":  T_NULL,
	"only":  T_ONLY, // used in media query
	"and":   T_LOGICAL_AND,
	"not":   T_LOGICAL_NOT,
	"or":    T_LOGICAL_OR,
	"xor":   T_LOGICAL_XOR,
}

var UnitTokenMap = KeywordTokenMap{
	"px":  T_UNIT_PX,
	"pt":  T_UNIT_PT,
	"pc":  T_UNIT_PC,
	"em":  T_UNIT_EM,
	"cm":  T_UNIT_CM,
	"ex":  T_UNIT_EX,
	"ch":  T_UNIT_CH,
	"in":  T_UNIT_IN,
	"mm":  T_UNIT_MM,
	"rem": T_UNIT_REM,
	"vh":  T_UNIT_VH,
	"vw":  T_UNIT_VW,

	"Hz":  T_UNIT_HZ,
	"kHz": T_UNIT_KHZ,

	"vmin": T_UNIT_VMIN,
	"vmax": T_UNIT_VMAX,
	"deg":  T_UNIT_DEG,
	"grad": T_UNIT_GRAD,
	"rad":  T_UNIT_RAD,
	"turn": T_UNIT_TURN,
	"dpi":  T_UNIT_DPI,
	"dpcm": T_UNIT_DPCM,
	"dppx": T_UNIT_DPPX,
	"s":    T_UNIT_SECOND,
	"ms":   T_UNIT_MILLISECOND,
	"%":    T_UNIT_PERCENT,
}

type Token struct {
	Type                  TokenType
	Str                   string
	Pos                   int
	Line                  int
	ContainsInterpolation bool
}

type TokenChannel chan *Token

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
		T_CLASS_SELECTOR, T_PARENT_SELECTOR, T_PSEUDO_SELECTOR,
		T_ADJACENT_SIBLING_COMBINATOR, T_GENERAL_SIBLING_COMBINATOR,
		T_CHILD_COMBINATOR, T_DESCENDANT_COMBINATOR,
		T_BRACKET_LEFT: // '[' is the first token of attribute selector.
		return true
	}
	return false
}

func (tok Token) IsFlagKeyword() bool {
	switch tok.Type {
	case T_DEFAULT, T_OPTIONAL, T_GLOBAL, T_IMPORTANT:
		return true
	}
	return false
}

func (tok Token) IsLengthUnit() bool {
	switch tok.Type {
	case T_UNIT_EM, T_UNIT_EX, T_UNIT_CH, T_UNIT_REM, T_UNIT_CM, T_UNIT_IN,
		T_UNIT_MM, T_UNIT_PC, T_UNIT_PT, T_UNIT_PX, T_UNIT_VH, T_UNIT_VW,
		T_UNIT_VMIN, T_UNIT_VMAX:
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
		T_UNIT_DPPX, T_UNIT_DEG, T_UNIT_GRAD, T_UNIT_RAD, T_UNIT_TURN:
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
	T_AND_SELECTOR                // {parent-selector}{child-selector} { }
	T_DESCENDANT_COMBINATOR       // E ' ' F
	T_CHILD_COMBINATOR            // E '>' F
	T_ADJACENT_SIBLING_COMBINATOR // E '+' F
	T_GENERAL_SIBLING_COMBINATOR  // E '~' F

	T_UNICODE_RANGE

	T_IF
	T_ELSE
	T_ELSE_IF // @else if
	T_INCLUDE // for @include
	T_EACH    // for @each
	T_WHEN    // for @when
	T_MIXIN
	T_FUNCTION
	T_FOR
	T_FOR_FROM
	T_FOR_THROUGH
	T_FOR_TO
	T_FOR_IN
	T_WHILE
	T_RETURN
	T_RANGE   // for '..'
	T_CONTENT // for '@content'

	// Flag token types
	T_GLOBAL
	T_DEFAULT
	T_IMPORTANT
	T_OPTIONAL

	T_FONT_FACE

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
	T_BRACE_START
	T_BRACE_END
	T_LANG_CODE // 'en', 'fr', 'fr-ca'
	T_BRACKET_LEFT
	T_ATTRIBUTE_NAME
	T_BRACKET_RIGHT

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

	T_VARIABLE

	T_IMPORT
	T_AT_RULE

	T_CHARSET
	T_QQ_STRING
	T_Q_STRING
	T_UNQUOTE_STRING
	T_PAREN_START
	T_PAREN_END
	T_CONSTANT
	T_INTEGER
	T_FLOAT

	T_CDO // for <!--
	T_CDC // for -->

	/*
		unit tokens
	*/
	T_UNIT_NONE
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
