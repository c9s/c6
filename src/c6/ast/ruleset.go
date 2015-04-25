package ast

type RuleSet struct {
	Selectors []Selector
	Block     DeclarationBlock
}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

func (self *RuleSet) AppendSelector(sel Selector) {
	self.Selectors = append(self.Selectors, sel)
}

func (self *RuleSet) AppendDeclaration() {

}

func (self *RuleSet) AppendSubRuleSet() {

}

// Complete the statement interface
func (self *RuleSet) IsStatement() {}

type PropertyName struct{}
type PropertyValue struct{}

type Property struct {
	Name   PropertyName
	Values []PropertyValue
}

func (self Property) addValue(value PropertyValue) {
	self.Values = append(self.Values, value)
}

/*
A declaration can be a property or a ruleset
*/
type Declaration interface{}

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	Declarations []Declaration
}

func (self DeclarationBlock) String() string {
	/*
		for _, property := range self.Properties {

		}
	*/
	return ""
}
