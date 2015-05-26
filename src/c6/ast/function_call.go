package ast

import "sort"

// type FunctionCallArgument Expression
type FunctionCallArgument struct {
	Argument                Expression
	ArgumentDefineReference *Argument
}

func (arg FunctionCallArgument) String() (out string) {
	out = arg.Argument.String()
	if arg.ArgumentDefineReference != nil {
		out += " ref:" + arg.ArgumentDefineReference.String()
	}
	return out
}

func NewFunctionCallArgument(expr Expression) *FunctionCallArgument {
	return &FunctionCallArgument{expr, nil}
}

type FunctionCallArguments []*FunctionCallArgument

func (args *FunctionCallArguments) Sort() {
	var sorter = &FunctionCallArgumentsSorter{*args}
	sort.Sort(sorter)
}

type FunctionCallArgumentsSorter struct {
	Arguments FunctionCallArguments
}

// Len is part of sort.Interface.
func (s *FunctionCallArgumentsSorter) Len() int {
	return len(s.Arguments)
}

// Swap is part of sort.Interface.
func (s *FunctionCallArgumentsSorter) Swap(i, j int) {
	s.Arguments[i], s.Arguments[j] = s.Arguments[j], s.Arguments[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *FunctionCallArgumentsSorter) Less(i, j int) bool {
	return s.Arguments[i].ArgumentDefineReference.Position < s.Arguments[j].ArgumentDefineReference.Position
	// return s.by(&s.Arguments[i], &s.Arguments[j])
}

type FunctionCall struct {
	Ident     *Token
	Arguments FunctionCallArguments
}

func (self FunctionCall) CanBeNode() {}
func (self FunctionCall) String() (out string) {
	out = self.Ident.Str + "("
	for _, arg := range self.Arguments {
		out += arg.Argument.String() + ", "
	}
	if len(self.Arguments) > 0 {
		out = out[:len(out)-2]
	}
	out += ")"
	return out
}

func (self *FunctionCall) AppendArgument(arg Expression) {
	self.Arguments = append(self.Arguments, NewFunctionCallArgument(arg))
}

func NewFunctionCallWithToken(token *Token) *FunctionCall {
	return &FunctionCall{
		Ident:     token,
		Arguments: FunctionCallArguments{},
	}
}
