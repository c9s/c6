package runtime

import "c6/ast"

type Context struct {
	RuleSetStack   []*ast.RuleSet
	GlobalSymTable SymTable
}

func (context *Context) GetVariable(name string) *ast.Variable {
	var idx = len(context.RuleSetStack) - 1
	for ; idx > 0; idx-- {
		// context.RuleSetStack
	}

	return nil
}

func (context *Context) TopRule() *ast.RuleSet {
	return context.RuleSetStack[len(context.RuleSetStack)-1]
}
