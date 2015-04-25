package ast

type Variable struct {
	Name      string
	Value     interface{}
	ScopeRule *RuleSet
}

func NewVariable(name string) *Variable {
	return &Variable{name, nil, nil}
}

func (self *Variable) SetScopeRule(scope *RuleSet) {
	self.ScopeRule = scope
}

func (self *Variable) SetValue(val interface{}) {
	self.Value = val
}
