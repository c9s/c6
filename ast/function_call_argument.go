package ast

import "sort"

// type FunctionCallArgument Expr
type FunctionCallArgument struct {
	Argument                Expr
	ArgumentDefineReference *Argument
}

func (arg FunctionCallArgument) String() (out string) {
	out = arg.Argument.String()
	if arg.ArgumentDefineReference != nil {
		out += " ref:" + arg.ArgumentDefineReference.String()
	}
	return out
}

func NewFunctionCallArgument(expr Expr) *FunctionCallArgument {
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
