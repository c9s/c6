package ast

/*
This is a SCSS only selector
*/
type ParentSelector struct {
	ParentRuleSet *RuleSet
	Token         *Token
}

func (self ParentSelector) String() string {
	// TODO: get parent rule set and render the selector...
	return "ParentSelector.String()"
}

func NewParentSelectorWithToken(parentRuleSet *RuleSet, token *Token) *ParentSelector {
	return &ParentSelector{parentRuleSet, token}
}
