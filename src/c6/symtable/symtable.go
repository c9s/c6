package symtable

/*
SymTableItem can be anything

This package doesn't use ast.* types because we have to avoid package acyclic
reference.
*/
type SymTableItem interface{}

type SymTable map[string]SymTableItem

func NewSymTable() *SymTable {
	return &SymTable{}
}

func (self SymTable) Set(name string, v SymTableItem) {
	self[name] = v
}

func (self SymTable) Get(name string) (SymTableItem, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self SymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
