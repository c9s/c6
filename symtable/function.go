package symtable

import "github.com/c9s/c6/ast"

/*
FunctionSymTableItem can be anything

This package doesn't use ast.* types because we have to avoid package acyclic
reference.
*/
type FunctionSymTable map[string]*ast.Function

// type FunctionFunctionSymTable map[string]*ast.Function

func NewFunctionSymTable() *FunctionSymTable {
	return &FunctionSymTable{}
}

func (self FunctionSymTable) Set(name string, v *ast.Function) {
	self[name] = v
}

func (self FunctionSymTable) Get(name string) (*ast.Function, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self FunctionSymTable) Merge(a *FunctionSymTable) {
	for key, val := range *a {
		self[key] = val
	}
}

func (self FunctionSymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
