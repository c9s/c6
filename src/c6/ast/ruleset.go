package ast

type RuleSet struct {
	Selectors *CombinedSelector
	Block     *DeclarationBlock
}

type PropertyName struct{}
type PropertyValue struct{}

type Property struct {
	Name   PropertyName
	Values []PropertyValue
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
