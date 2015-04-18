package c6

type TokenType int

type Token struct {
	Type TokenType
	Str  string
	Pos  int
	Line int
}

const (
	T_SPACE = iota
	T_COMMENT_LINE
	T_COMMENT_BLOCK
	T_ID_SELECTOR
	T_CLASS_SELECTOR
	T_TAGNAME_SELECTOR
	T_VARIABLE
	T_AT_IMPORT
	T_QQ_STRING
	T_Q_STRING
)
