package c6

//go:generate stringer -type=TokenType
type TokenType int

const LF = '\r'
const CR = '\n'

type Token struct {
	Type TokenType
	Str  string
	Pos  int
	Line int
}

const (
	T_SPACE TokenType = iota
	T_COMMENT_LINE
	T_COMMENT_BLOCK
	T_SEMICOLON
	T_COMMA
	T_ID_SELECTOR
	T_CLASS_SELECTOR
	T_TAGNAME_SELECTOR
	T_PARENT_SELECTOR // SASS parent selector
	T_BRACE_START
	T_BRACE_END
	T_VARIABLE
	T_IMPORT
	T_CHARSET
	T_QQ_STRING
	T_Q_STRING
	T_PAREN_START
	T_PAREN_END
	T_CONSTANT
	T_INTEGER
	T_FLOAT
	T_UNIT_PX
	T_UNIT_PT
	T_UNIT_EM
	T_PROPERTY_NAME
	T_PROPERTY_VALUE
	T_HEX_COLOR
	T_COLON
)
