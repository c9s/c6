package ast

type Interpolation struct {
	Expr       Expr
	StartToken *Token
	EndToken   *Token
}

func (self Interpolation) CanBeNode() {}

func (self Interpolation) String() string {
	return self.Expr.String()
}

func NewInterpolation(expr Expr, startToken *Token, endToken *Token) *Interpolation {
	return &Interpolation{expr, startToken, endToken}
}
