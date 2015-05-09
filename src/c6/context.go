package c6

import "c6/symtable"
import "c6/ast"

/**
The Context contains all runtime variables and ruleset stack
*/
type Context struct {
	RuleSetStack   []*ast.RuleSet
	GlobalSymTable *symtable.SymTable
}

func NewContext() *Context {
	return &Context{[]*ast.RuleSet{}, &symtable.SymTable{}}
}

func (context *Context) PushRuleSet(ruleSet *ast.RuleSet) {
	var newStack = append(context.RuleSetStack, ruleSet)
	context.RuleSetStack = newStack
}

func (context *Context) PopRuleSet() (*ast.RuleSet, bool) {
	if len(context.RuleSetStack) == 0 {
		// XXX: throw error here?
		return nil, false
	} else if len(context.RuleSetStack) == 1 {
		ruleset := context.RuleSetStack[0]

		// clear the ruleset
		context.RuleSetStack = []*ast.RuleSet{}
		return ruleset, true
	}
	var idx = len(context.RuleSetStack) - 1
	var ruleSet = context.RuleSetStack[idx]
	context.RuleSetStack = context.RuleSetStack[:idx-1]
	return ruleSet, true
}

func (context *Context) GetVariable(name string) (symtable.SymTableItem, bool) {
	var idx = len(context.RuleSetStack) - 1
	for ; idx > 0; idx-- {
		ruleset := context.RuleSetStack[idx]
		if variable, ok := ruleset.Block.SymTable.Get(name); ok {
			return variable, true
		}
	}
	return nil, false
}

func (context *Context) TopRuleSet() *ast.RuleSet {
	if len(context.RuleSetStack) > 0 {
		return context.RuleSetStack[len(context.RuleSetStack)-1]
	}
	return nil
}
