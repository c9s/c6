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
