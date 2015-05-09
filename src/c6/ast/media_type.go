package ast

type MediaType struct {
	Name  string
	Token *Token
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
