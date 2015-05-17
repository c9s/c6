package ast

/**
@see http://www.w3.org/TR/CSS21/grammar.html


UniversalSelector
TypeSelector
DescendantCombinator
PseudoSelector
ChildCombinator
ClassSelector
IdSelector
AdjacentCombinator
AttributeSelector

*/
type Selector interface {
	String() string
}

type Combinator interface{}

// one or more simple selector
type CompoundSelector []Selector

func (self *CompoundSelector) Append(sel Selector) {
	var slice = append(*self, sel)
	*self = slice
}

func (self CompoundSelector) Length() int {
	return len(self)
}

func (self CompoundSelector) String() string {
	return "CompoundSelector.String()"
}

func NewCompoundSelector() *CompoundSelector {
	return &CompoundSelector{}
}

type ComplexSelectorItem struct {
	Combinator       Combinator
	CompoundSelector Selector
}

type ComplexSelector struct {
	CompoundSelector     Selector
	ComplexSelectorItems []*ComplexSelectorItem
}

func (self *ComplexSelector) AppendCompoundSelector(comb Combinator, sel Selector) {
	self.ComplexSelectorItems = append(self.ComplexSelectorItems, &ComplexSelectorItem{comb, sel})
}

func NewComplexSelector(sel Selector) *ComplexSelector {
	return &ComplexSelector{
		CompoundSelector: sel,
	}
}

/*
SelectorList struct
*/
type SelectorList []Selector

func (self *SelectorList) Append(sel Selector) {
	var slice = append(*self, sel)
	*self = slice
}

func (self SelectorList) Length() int {
	return len(self)
}

func (self SelectorList) String() (out string) {
	for _, sel := range self {
		out += sel.String()
	}
	return out
}

func NewSelectorList() *SelectorList {
	return &SelectorList{}
}

/**
TypeSelector
*/
type TypeSelector struct {
	Type  string
	Token *Token
}

func (self TypeSelector) String() string {
	return self.Type
}

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

func (self IdSelector) String() string {
	return self.Id
}

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

func (self ClassSelector) String() string {
	return self.ClassName
}

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
		return "[" + self.Name.String() + self.Match.String() + self.Pattern.String() + "]"
	}
	return "[" + self.Name.String() + "]"
}

func NewAttributeSelector(name, match, pattern *Token) *AttributeSelector {
	return &AttributeSelector{name, match, pattern}
}

func NewAttributeSelectorNameOnly(name *Token) *AttributeSelector {
	return &AttributeSelector{name, nil, nil}
}

type UniversalSelector struct {
	Token *Token
}

func (self UniversalSelector) String() string {
	return "*"
}

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
	C           string // for dynamic language pseudo selector like :lang(C)
	Token       *Token
}

func (self PseudoSelector) String() (out string) {
	if self.C != "" {
		return ":" + self.PseudoClass + "(" + self.C + ")"
	}
	return ":" + self.PseudoClass
}

func NewPseudoSelectorWithToken(token *Token) *PseudoSelector {
	return &PseudoSelector{token.Str, "", token}
}

/*
Selectors present: E '+' F
*/
type AdjacentCombinator struct {
	Token *Token
}

func (self AdjacentCombinator) String() string { return " + " }

func NewAdjacentCombinatorWithToken(token *Token) *AdjacentCombinator {
	return &AdjacentCombinator{token}
}

type DescendantCombinator struct {
	Token *Token
}

func (self DescendantCombinator) String() string { return " " }

func NewDescendantCombinatorWithToken(token *Token) *DescendantCombinator {
	return &DescendantCombinator{token}
}

func NewDescendantCombinator() *DescendantCombinator {
	return &DescendantCombinator{}
}

type GeneralSiblingCombinator struct{ Token *Token }

func NewGeneralSiblingCombinator() *GeneralSiblingCombinator {
	return &GeneralSiblingCombinator{}
}

func NewGeneralSiblingCombinatorWithToken(token *Token) *GeneralSiblingCombinator {
	return &GeneralSiblingCombinator{token}
}

func (self GeneralSiblingCombinator) String() string { return " ~ " }

type ChildCombinator struct {
	Token *Token
}

func (self ChildCombinator) String() string { return " > " }

func NewChildCombinatorWithToken(token *Token) *ChildCombinator {
	return &ChildCombinator{token}
}

func NewChildCombinator() *ChildCombinator {
	return &ChildCombinator{}
}

/*
This is a SCSS only selector
*/
type ParentSelector struct {
	ParentRuleSet *RuleSet
	Token         *Token
}

func (self ParentSelector) String() string {
	// TODO: get parent rule set and render the selector...
	panic("unimplemented")
	return ""
}

func NewParentSelectorWithToken(parentRuleSet *RuleSet, token *Token) *ParentSelector {
	return &ParentSelector{parentRuleSet, token}
}
