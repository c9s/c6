package runtime

import "github.com/c9s/c6/symtable"
import "github.com/c9s/c6/ast"

/**
The Context contains all runtime variables and ruleset stack
*/
type Context struct {
	RuleSetStack []*ast.RuleSet
	GlobalBlock  *ast.Block
	Mixins       *symtable.SymTable
	Functions    *symtable.SymTable
	Stack        *Stack
	ImportedPath map[string]bool
}

func NewContext() *Context {
	return &Context{
		GlobalBlock: ast.NewBlock(),
		Mixins:      &symtable.SymTable{},
		Functions:   &symtable.SymTable{},
		Stack:       &Stack{},
	}
}

func (context *Context) CurrentBlock() ast.BlockNode {
	if ruleSet := context.TopRuleSet(); ruleSet != nil {
		return ruleSet.Block
	}
	return context.GlobalBlock
}

func (context *Context) PushRuleSet(ruleSet *ast.RuleSet) {
	var newStack = append(context.RuleSetStack, ruleSet)
	context.RuleSetStack = newStack
}

func (context *Context) PopRuleSet() (*ast.RuleSet, bool) {
	if len(context.RuleSetStack) == 0 {
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

func (context *Context) TopRuleSet() *ast.RuleSet {
	if len(context.RuleSetStack) > 0 {
		return context.RuleSetStack[len(context.RuleSetStack)-1]
	}
	return nil
}

func (context *Context) GetVariable(name string) (interface{}, bool) {
	var idx = len(context.RuleSetStack) - 1
	for ; idx > 0; idx-- {
		/*
			FIXME:
			ruleset := context.RuleSetStack[idx]
			if variable, ok := ruleset.Block.SymTable.Get(name); ok {
				return variable, true
			}
		*/
	}
	/*
		FIXME:
		if variable, ok := context.GlobalBlock.SymTable.Get(name); ok {
			return variable, true
		}
	*/
	return nil, false
}
