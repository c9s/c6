package symtable

/*
SymTableItem can be anything

This package doesn't use ast.* types because we have to avoid package acyclic
reference.
*/
type SymTable map[string]interface{}

// type FunctionSymTable map[string]*ast.Function

func NewSymTable() *SymTable {
	return &SymTable{}
}

func (self SymTable) Set(name string, v interface{}) {
	self[name] = v
}

func (self SymTable) Get(name string) (interface{}, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self SymTable) Merge(a *SymTable) {
	for key, val := range *a {
		self[key] = val
	}
}

func (self SymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
