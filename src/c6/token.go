package c6

//go:generate stringer -type=TokenType token.go
type TokenType int

const LF = '\r'
const CR = '\n'

type Token struct {
	Type                  TokenType
	Str                   string
	Pos                   int
	Line                  int
	ContainsInterpolation bool
}

func (tok Token) IsString() bool {
	return tok.Type == T_QQ_STRING || tok.Type == T_Q_STRING || tok.Type == T_UNQUOTE_STRING
}

func (tok Token) IsSelectorOp() bool {
	return tok.Type == T_ADJACENT_SELECTOR ||
		tok.Type == T_CHILD_SELECTOR ||
		tok.Type == T_DESCENDANT_SELECTOR
}

func (tok Token) IsSelector() bool {
	return tok.Type == T_TYPE_SELECTOR ||
		tok.Type == T_UNIVERSAL_SELECTOR ||
		tok.Type == T_ID_SELECTOR ||
		tok.Type == T_CLASS_SELECTOR ||
		tok.Type == T_PARENT_SELECTOR ||
		tok.Type == T_PSEUDO_SELECTOR ||
		tok.Type == T_ADJACENT_SELECTOR ||
		tok.Type == T_CHILD_SELECTOR ||
		tok.Type == T_DESCENDANT_SELECTOR ||
		tok.Type == T_PARENT_SELECTOR
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

	// selector tokens
	T_ID_SELECTOR
	T_CLASS_SELECTOR
	T_TYPE_SELECTOR
	T_UNIVERSAL_SELECTOR
	T_PARENT_SELECTOR        // SASS parent selector
	T_PSEUDO_SELECTOR        // :hover, :visited , ...
	T_INTERPOLATION_SELECTOR // selector with interpolation: '#{ ... }'
	T_CONCAT                 // used to concat selectors and interpolation

	// Selector relationship
	T_AND_SELECTOR        // {parent-selector}{child-selector} { }
	T_DESCENDANT_SELECTOR // E ' ' F
	T_CHILD_SELECTOR      // E '>' F
	T_ADJACENT_SELECTOR   // E '+' F

	T_OR   // 'or' used in conditional query
	T_AND  // 'and' used in conditional query
	T_PLUS // E '+' F
	T_GT   // E '>' F
	T_BRACE_START
	T_BRACE_END
	T_LANG_CODE // 'en', 'fr', 'fr-ca'
	T_BRACKET_LEFT
	T_ATTRIBUTE_NAME
	T_BRACKET_RIGHT
	T_EQUAL       // for '='
	T_TILDE_EQUAL // for '~='
	T_PIPE_EQUAL  // for '|='
	T_VARIABLE
	T_IMPORT
	T_CHARSET
	T_QQ_STRING
	T_Q_STRING
	T_UNQUOTE_STRING
	T_PAREN_START
	T_PAREN_END
	T_CONSTANT
	T_INTEGER
	T_FLOAT
	T_UNIT_PX
	T_UNIT_PT
	T_UNIT_EM
	T_UNIT_REM
	T_UNIT_DEG
	T_UNIT_PERCENT
	T_PROPERTY_NAME
	T_PROPERTY_VALUE
	T_HEX_COLOR
	T_COLON
	T_INTERPOLATION_START
	T_INTERPOLATION_INNER
	T_INTERPOLATION_END
	T_DIV
	T_MUL
	T_MINUS
)
