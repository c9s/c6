package c6

import "c6/ast"
import "container/list"

/**
The Context contains all runtime variables and ruleset stack
*/
type Context struct {
	RuleSetList []*ast.RuleSet

	// SymTableStack  []*ast.SymTable
	GlobalSymTable ast.SymTable
}

func NewContext() *Context {
	var ruleSetStack = list.New()
	var context = &Context{ruleSetStack, ast.SymTable{}}
	return context
}

func (context *Context) PushRuleSet(ruleSet *ast.RuleSet) {
	var newStack = append(context.RuleSetList, ruleSet)
	context.RuleSetStack = newStack
}

func (context *Context) PopRuleSet() *ast.RuleSet {
	if len(context.RuleSetStack) == 0 {
		// XXX: throw error here?
		return nil
	}
	var idx = len(context.RuleSetStack) - 1
	ruleSet := context.RuleSetStack[idx]
	context.RuleSetStack = context.RuleSetStack[:idx-1]
	return ruleSet
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
