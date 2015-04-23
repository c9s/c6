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

const ChildSelectorOp = " > "
const DescendantSelectorOp = " "
const AdjacentSelectorOp = " + "

type Selector interface {
	String() string
}

type CodeGen interface{}

type UniversalSelector struct{}

func (self UniversalSelector) String() string {
	return "*"
}

type PseudoSelector struct {
	PseudoClass string
	C           string
}

func (self PseudoSelector) String() (out string) {
	if self.C != "" {
		return ":" + self.PseudoClass + "(" + self.C + ")"
	}
	return ":" + self.PseudoClass
}

/**
TypeSelector
*/
type TypeSelector struct {
	Type string
}

func (self TypeSelector) String() string {
	return self.Type
}

type IdSelector struct {
	Id string
}

func (self IdSelector) String() string {
	return "#" + self.Id
}

type ClassSelector struct {
	ClassName string
}

func (self ClassSelector) String() string {
	return "." + self.ClassName
}

type AttributeSelector struct {
	Name    string
	Op      string
	Pattern string
}

func (self AttributeSelector) String() (out string) {
	if self.Op != "" && self.Pattern != "" {
		return "[" + self.Name + self.Op + self.Pattern + "]"
	}
	return "[" + self.Name + "]"
}

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
