package ast

type MediaType struct {
	Expr Expression
}

func NewMediaType(expr Expression) *MediaType {
	return &MediaType{expr}
}

func (self MediaType) String() string {
	return self.Expr.String()
}

type MediaFeature struct {
	Feature Expression
	Value   Expression
	Open    *Token
	Close   *Token
}

func NewMediaFeature(feature, value Expression) *MediaFeature {
	return &MediaFeature{Feature: feature, Value: value}
}

func (self MediaFeature) String() (out string) {
	out = "(" + self.Feature.String()
	if self.Value != nil {
		out += ":" + self.Value.String()
	}
	out += ")"
	return out
}

/*
  media_type: all | aural | braille | handheld | print |
  projection | screen | tty | tv | embossed
*/
const (
	MediaTypeAll = iota
	MediaTypeAural
	MediaTypeBraille
	MediaTypeHandheld
	MediaTypePrint
	MediaTypeProjection
	MediaTypeScreen
	MediaTypeTTY
	MediaTypeTV
	MediaTypeEmbossed
)
