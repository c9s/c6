package runtime

import "c6/ast"

type SymTable map[string]*ast.Variable

func (self SymTable) AddVariable(v *ast.Variable) {
	self[v.Name] = v
}

func (self SymTable) FindVariable(name string) *ast.Variable {
	if val, ok := self[name]; ok {
		return val
	}
	return nil
}

func (self SymTable) HasVariable(v *ast.Variable) bool {
	if _, ok := self[v.Name]; ok {
		return true
	} else {
		return false
	}
}
