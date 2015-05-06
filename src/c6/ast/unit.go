package ast

type Unit struct {
	Type  TokenType
	Token *Token
}

func NewUnitWithToken(token *Token) *Unit {
	return &Unit{token.Type, token}
}

func (unit Unit) String() string {
	if unit.Token != nil {
		return unit.Token.Str
	}
	return string(unit.Type.String())
}
