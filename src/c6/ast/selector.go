package ast

/**
@see http://www.w3.org/TR/CSS21/grammar.html

UniversalSelector
TypeSelector
PseudoSelector
ClassSelector
IdSelector
AttributeSelector
*/
type Selector interface {
	String() string
}

/**
TypeSelector
*/
type TypeSelector struct {
	Type  string
	Token *Token
}

func (self TypeSelector) String() string    { return self.Type }
func (self TypeSelector) CSSString() string { return self.Type }

func NewTypeSelectorWithToken(token *Token) *TypeSelector {
	return &TypeSelector{token.Str, token}
}

func NewTypeSelector(typename string) *TypeSelector {
	return &TypeSelector{typename, nil}
}

type IdSelector struct {
	Id    string
	Token *Token
}

func (self IdSelector) String() string    { return self.Id }
func (self IdSelector) CSSString() string { return self.Id }

func NewIdSelectorWithToken(token *Token) *IdSelector {
	return &IdSelector{token.Str, token}
}

func NewIdSelector(id string) *IdSelector {
	return &IdSelector{id, nil}
}

type ClassSelector struct {
	ClassName string
	Token     *Token
}

func (self ClassSelector) String() string    { return self.ClassName }
func (self ClassSelector) CSSString() string { return self.ClassName }

func NewClassSelectorWithToken(token *Token) *ClassSelector {
	return &ClassSelector{token.Str, token}
}

func NewClassSelector(className string) *ClassSelector {
	return &ClassSelector{className, nil}
}

type AttributeSelector struct {
	Name    *Token
	Match   *Token
	Pattern *Token
}

func (self AttributeSelector) String() (out string) {
	if self.Match != nil && self.Pattern != nil {
		return "[" + self.Name.Str + self.Match.Str + self.Pattern.Str + "]"
	}
	return "[" + self.Name.Str + "]"
}

func (self AttributeSelector) CSSString() (out string) { return self.String() }

func NewAttributeSelector(name, match, pattern *Token) *AttributeSelector {
	return &AttributeSelector{name, match, pattern}
}

func NewAttributeSelectorNameOnly(name *Token) *AttributeSelector {
	return &AttributeSelector{name, nil, nil}
}

type UniversalSelector struct {
	Token *Token
}

func (self UniversalSelector) String() string    { return "*" }
func (self UniversalSelector) CSSString() string { return "*" }

func NewUniversalSelectorWithToken(token *Token) *UniversalSelector {
	return &UniversalSelector{token}
}

func NewUniversalSelector() *UniversalSelector {
	return &UniversalSelector{}
}

/*
Selectors presents: E:pseudo
*/
type PseudoSelector struct {
	PseudoClass string
	Token       *Token
}

func (self PseudoSelector) String() (out string)    { return ":" + self.PseudoClass }
func (self PseudoSelector) CSSString() (out string) { return ":" + self.PseudoClass }

func NewPseudoSelectorWithToken(token *Token) *PseudoSelector {
	return &PseudoSelector{token.Str, token}
}

/*
Selectors presents: E:pseudo
*/
type FunctionalPseudoSelector struct {
	PseudoClass string
	C           string // for dynamic language pseudo selector like :lang(C)
	Token       *Token
}

func (self FunctionalPseudoSelector) String() (out string) {
	if self.C != "" {
		return ":" + self.PseudoClass + "(" + self.C + ")"
	}
	return ":" + self.PseudoClass
}

func (self FunctionalPseudoSelector) CSSString() string { return self.String() }

func NewFunctionalPseudoSelectorWithToken(token *Token) *FunctionalPseudoSelector {
	return &FunctionalPseudoSelector{token.Str, "", token}
}
