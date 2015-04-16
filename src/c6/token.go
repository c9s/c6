package c6

type Token struct {
	Type int
	Str  string
	Pos  int
	Line int
}

const (
	TokenClassSelector = iota
	TokenSpace
)
