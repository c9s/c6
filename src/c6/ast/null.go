package ast

/*
Null struct presents Null value
*/
type Null struct {
	Token *Token
}

func (self Null) String() string {
	return "null"
}

func NewNullWithToken(tok *Token) *Null {
	return &Null{tok}
}
