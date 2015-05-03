package c6

import "c6/ast"

type Context struct {
	RuleSetStack   []*ast.RuleSet
	SymTableStack  []*ast.SymTable
	GlobalSymTable ast.SymTable
}

func NewContext() *Context {
	return &Context{[]*ast.RuleSet{}, []*ast.SymTable{}, ast.SymTable{}}
}

func (context *Context) GetVariable(name string) *ast.Variable {
	var idx = len(context.SymTableStack) - 1
	for ; idx > 0; idx-- {
		stack := context.SymTableStack[idx]
		if variable := stack.FindVariable(name); variable != nil {
			return variable
		}
	}
	return nil
}

func (context *Context) TopRuleSet() *ast.RuleSet {
	if len(context.RuleSetStack) > 0 {
		return context.RuleSetStack[len(context.RuleSetStack)-1]
	}
	return nil
}
