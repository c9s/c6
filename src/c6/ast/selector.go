package ast

import "strings"

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

type UniversalSelector struct {
	Token *Token
}

func (self UniversalSelector) IsSelector() {}

func (self UniversalSelector) String() string {
	return "*"
}

func NewUniversalSelector(token *Token) *UniversalSelector {
	return &UniversalSelector{token}
}

type DescendantCombinator struct{}

func (self DescendantCombinator) IsSelector()    {}
func (self DescendantCombinator) String() string { return " " }

func NewDescendantCombinator() *DescendantCombinator {
	return &DescendantCombinator{}
}

type ChildCombinator struct{}

func (self ChildCombinator) IsSelector()    {}
func (self ChildCombinator) String() string { return " > " }

func NewChildCombinator() *ChildCombinator {
	return &ChildCombinator{}
}

/*
Selectors presents: E:pseudo
*/
type PseudoSelector struct {
	PseudoClass string
	C           string // for dynamic language pseudo selector like :lang(C)
}

func (self PseudoSelector) IsSelector() {}
func (self PseudoSelector) String() (out string) {
	if self.C != "" {
		return ":" + self.PseudoClass + "(" + self.C + ")"
	}
	return ":" + self.PseudoClass
}

/*
Selectors present: E '+' F
*/
type AdjacentCombinator struct{}

func (self AdjacentCombinator) IsSelector()    {}
func (self AdjacentCombinator) String() string { return " + " }

/**
TypeSelector
*/
type TypeSelector struct {
	Type string
}

func (self TypeSelector) IsSelector() {}
func (self TypeSelector) String() string {
	return self.Type
}

type IdSelector struct {
	Id string
}

func (self IdSelector) IsSelector() {}
func (self IdSelector) String() string {
	return self.Id
}

type ClassSelector struct {
	ClassName string
}

func (self ClassSelector) IsSelector() {}
func (self ClassSelector) String() string {
	return self.ClassName
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

/*
This is a SCSS only selector
*/
type ParentSelector struct {
	ParentRuleSet *RuleSet
}

func (self ParentSelector) IsSelector() {}
func (self ParentSelector) String() string {
	// TODO: get parent rule set and render the selector...
	panic("unimplemented")
	return ""
}

/**
An ast node that could combine all selector with the same operator.
*/
type CombinedSelector struct {
	Op        string
	Selectors []Selector
}

func (self *CombinedSelector) addSelector(sel Selector) {
	self.Selectors = append(self.Selectors, sel)
}

func (self CombinedSelector) String() string {
	var out []string = []string{}
	for _, sel := range self.Selectors {
		out = append(out, sel.String())
	}
	return strings.Join(out, self.Op)
}
