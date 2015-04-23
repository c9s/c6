package ast

type RuleSet struct {
	Selectors *CombinedSelector
	Block     *DeclarationBlock
}

type Property struct{}

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}

*/
type DeclarationBlock struct {
	Properties []Property
	SubRules   []RuleSet
}

func (self DeclarationBlock) String() string {
	/*
		for _, property := range self.Properties {

		}
	*/
	return ""
}
