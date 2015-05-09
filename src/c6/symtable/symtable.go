package symtable

/*
SymTableItem can be anything
*/
type SymTableItem interface{}

type SymTable map[string]SymTableItem

func NewSymTable() *SymTable {
	return &SymTable{}
}

func (self SymTable) Set(name string, v SymTableItem) {
	self[name] = v
}

func (self SymTable) Get(name string) SymTableItem {
	if val, ok := self[name]; ok {
		return val
	}
	return nil
}

func (self SymTable) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
