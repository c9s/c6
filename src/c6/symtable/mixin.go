package symtable

import "c6/ast"

/*
This package doesn't use ast.* types because we have to avoid package acyclic
reference.
*/
type MixinSymTable map[string]*ast.MixinStatement

// type FunctionMixinSymTable map[string]*ast.MixinStatement

func NewMixinSymTable() *MixinSymTable {
	return &MixinSymTable{}
}

func (self MixinSymTable) Set(name string, v *ast.MixinStatement) {
	self[name] = v
}

func (self MixinSymTable) Get(name string) (*ast.MixinStatement, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self MixinSymTable) Merge(a *MixinSymTable) {
	for key, val := range *a {
		self[key] = val
	}
}

func (self MixinSymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
