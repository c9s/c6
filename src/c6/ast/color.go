package ast

type HexColor struct {
	Hex   string
	Token *Token
}

func (self HexColor) CanBeExpression() {}

// Factor functions
func NewHexColor(hex string, token *Token) *HexColor {
	return &HexColor{hex, token}
}

type RGBAColor struct {
	R     float64
	G     float64
	B     float64
	A     float64
	Token *Token
}

func (self RGBAColor) CanBeExpression() {}

func NewRGBAColor(r, g, b, a float64, token *Token) *RGBAColor {
	return &RGBAColor{r, g, b, a, token}
}
