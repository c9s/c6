package ast

import "strings"

/**
@see http://www.w3.org/TR/CSS21/grammar.html


UniversalSelector
TypeSelector
DescendantSelector
PseudoSelector
ChildSelector
ClassSelector
IdSelector
AdjacentSelector
AttributeSelector
*/
type Selector interface{}

type CodeGen interface{}

type UniversalSelector struct{}

func (self *UniversalSelector) String() string {
	return "*"
}

/**
TypeSelector
*/
type TypeSelector struct {
	Type string
}

func (self *TypeSelector) String() string {
	return self.Type
}

type IdSelector struct {
	Id string
}

func (self *IdSelector) String() string {
	return "#" + self.Id
}

type ClassSelector struct {
	ClassName string
}

func (self *ClassSelector) String() string {
	return "." + self.ClassName
}

type CombinedSelector struct {
	Op        string
	Selectors []Selector
}

func (self *CombinedSelector) addSelector(sel Selector) {
	self.Selectors = append(self.Selectors, sel)
}

func (self *CombinedSelector) String() string {
	var out []string = []string{}
	for _, sel := range self.Selectors {
		out = append(out, sel.String())
	}
	return strings.Join(self.Op, out)
}
