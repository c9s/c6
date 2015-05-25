package ast

type ArgumentList struct {
	Arguments []*Argument
	Keywords  map[string]*Argument
}

// type ArgumentList []*Argument

func (self ArgumentList) Add(arg *Argument) {
	self.Arguments = append(self.Arguments, arg)
	self.Keywords[arg.VariableName.Str] = arg
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{
		Arguments: []*Argument{},
		Keywords:  map[string]*Argument{},
	}
}
