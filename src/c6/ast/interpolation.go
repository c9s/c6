package ast

type Interpolation struct {
	Expression Expression
	StartToken *Token
	EndToken   *Token
}

func (self Interpolation) CanBeNode() {}

func (self Interpolation) String() string {
	return self.Expression.String()
}

func NewInterpolation(expr Expression, startToken *Token, endToken *Token) *Interpolation {
	return &Interpolation{expr, startToken, endToken}
}
