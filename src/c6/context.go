package c6

import "c6/ast"

type Context struct {
	RuleSetStack []*ast.RuleSet
	SymTable     SymTable
}

func (context *Context) TopRule() *ast.RuleSet {
	return context.RuleSetStack[len(context.RuleSetStack)-1]
}
