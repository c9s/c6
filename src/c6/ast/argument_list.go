package ast

type ArgumentList struct {
	Arguments []*Argument
	Keywords  map[string]*Argument
}

// type ArgumentList []*Argument

func (self *ArgumentList) Add(arg *Argument) {
	arg.Position = len(self.Arguments)
	self.Arguments = append(self.Arguments, arg)
	self.Keywords[arg.Name.Str] = arg
}

func (self *ArgumentList) Lookup(name string) *Argument {
	if arg, ok := self.Keywords[name]; ok {
		return arg
	}
	return nil
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{
		Arguments: []*Argument{},
		Keywords:  map[string]*Argument{},
	}
}
