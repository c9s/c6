package ast

type String struct {
	// Can be `"`, `'` or ``
	Quote byte
	Value string
	Token *Token
}


func (self String) String() string {
	return self.Value
}

func NewStringWithQuote(quote byte, token *Token) *String {
	return &String{quote, token.Str, token}
}

func NewString(token *Token) *String {
	return &String{0, token.Str, token}
}
