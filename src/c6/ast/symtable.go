package ast

/*
XXX: This smells bad, and but we don't want to put all ast node types in one classes...
*/
type SymTable map[string]*Variable

func (self SymTable) AddVariable(v *Variable) {
	self[v.Name] = v
}

func (self SymTable) FindVariable(name string) *Variable {
	if val, ok := self[name]; ok {
		return val
	}
	return nil
}

func (self SymTable) HasVariable(v *Variable) bool {
	if _, ok := self[v.Name]; ok {
		return true
	} else {
		return false
	}
}
