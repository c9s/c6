package ast

type MediaType struct {
	Name  string
	Token *Token
}

func NewMediaType(name string) *MediaType {
	return &MediaType{name, nil}
}

func NewMediaTypeWithToken(token *Token) *MediaType {
	return &MediaType{token.Str, token}
}

func (self MediaType) String() string {
	return self.Name
}

type MediaFeature struct {
	Feature Expression
	Value   Expression
	Token   *Token
}

func NewMediaFeature(feature, value Expression) *MediaFeature {
	return &MediaFeature{feature, value, nil}
}

func NewMediaFeatureWithToken(feature, value Expression, token *Token) *MediaFeature {
	return &MediaFeature{feature, value, token}
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
