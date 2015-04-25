package ast

type FunctionCallArgument struct {
}

type FunctionCall struct {
	Function string
}

func (self FunctionCall) String() string {
	return self.Function + "(" + ")"
}
