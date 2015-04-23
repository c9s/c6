package ast

type RuleSet struct {
	selectors SelectorGroup
	block     DeclarationBlock
}

type DeclarationBlock struct{}
