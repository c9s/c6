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
	// type signature method
	IsSelector()

	// basic code gen
	String() string
}

/**
TypeSelector
*/
type TypeSelector struct {
	Type  string
	Token *Token
}

func (self TypeSelector) IsSelector() {}
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

func (self IdSelector) IsSelector() {}
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

func (self ClassSelector) IsSelector() {}
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
	Name    string
	Op      string
	Pattern string
}

func (self AttributeSelector) IsSelector() {}
func (self AttributeSelector) String() (out string) {
	if self.Op != "" && self.Pattern != "" {
		return "[" + self.Name + self.Op + self.Pattern + "]"
	}
	return "[" + self.Name + "]"
}

type UniversalSelector struct {
	Token *Token
}

func (self UniversalSelector) IsSelector() {}

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

func (self PseudoSelector) IsSelector() {}
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

func (self AdjacentCombinator) IsSelector()    {}
func (self AdjacentCombinator) String() string { return " + " }

func NewAdjacentCombinatorWithToken(token *Token) *AdjacentCombinator {
	return &AdjacentCombinator{token}
}

type DescendantCombinator struct {
	Token *Token
}

func (self DescendantCombinator) IsSelector()    {}
func (self DescendantCombinator) String() string { return " " }

func NewDescendantCombinatorWithToken(token *Token) *DescendantCombinator {
	return &DescendantCombinator{token}
}

func NewDescendantCombinator() *DescendantCombinator {
	return &DescendantCombinator{}
}

type GroupCombinator struct {
	Token *Token
}

func (self GroupCombinator) IsSelector()    {}
func (self GroupCombinator) String() string { return ", " }

func NewGroupCombinatorWithToken(token *Token) *GroupCombinator {
	return &GroupCombinator{token}
}

func NewGroupCombinator() *GroupCombinator {
	return &GroupCombinator{}
}

type ChildCombinator struct {
	Token *Token
}

func (self ChildCombinator) IsSelector()    {}
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

func (self ParentSelector) IsSelector() {}
func (self ParentSelector) String() string {
	// TODO: get parent rule set and render the selector...
	panic("unimplemented")
	return ""
}

func NewParentSelectorWithToken(parentRuleSet *RuleSet, token *Token) *ParentSelector {
	return &ParentSelector{parentRuleSet, token}
}
