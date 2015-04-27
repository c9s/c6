package ast

type Variable struct {
	Name      string
	Value     Expression
	ScopeRule *RuleSet
	Token     *Token
}

func (self Variable) CanBeExpression() {}
func (self Variable) CanBeNode()       {}

func (self *Variable) SetScopeRule(scope *RuleSet) {
	self.ScopeRule = scope
}

func (self *Variable) SetValue(val Expression) {
	self.Value = val
}

func NewVariable(token *Token) *Variable {
	return &Variable{token.Str, nil, nil, token}
}
