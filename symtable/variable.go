package symtable

import "github.com/c9s/c6/ast"

/*
This package doesn't use ast.* types because we have to avoid package acyclic
reference.
*/
type VarSymTable map[string]*ast.Expr

// type FunctionVarSymTable map[string]*ast.Expr

func NewVarSymTable() *VarSymTable {
	return &VarSymTable{}
}

func (self VarSymTable) Set(name string, v *ast.Expr) {
	self[name] = v
}

func (self VarSymTable) Get(name string) (*ast.Expr, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self VarSymTable) Merge(a *VarSymTable) {
	for key, val := range *a {
		self[key] = val
	}
}

func (self VarSymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
