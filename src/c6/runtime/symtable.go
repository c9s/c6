package runtime

import "c6/ast"

type SymTable map[string]interface{}

func (self SymTable) AddVariable(v *ast.Variable) {
	self[v.Name] = v
}

func (self SymTable) HasVariable(v *ast.Variable) bool {
	if _, ok := self[v.Name]; ok {
		return true
	} else {
		return false
	}
}
