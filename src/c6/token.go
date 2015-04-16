package c6

type TokenType int

type Token struct {
	Type TokenType
	Str  string
	Pos  int
	Line int
}

const (
	T_CLASS_SELECTOR = iota
	T_SPACE
)
