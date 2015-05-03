package ast

type FunctionCall struct {
	Function  string
	Arguments []Expression
	Token     *Token
}

func (self FunctionCall) CanBeNode()       {}
func (self FunctionCall) String() (out string) {
	out = self.Function + "("
	for _, arg := range self.Arguments {
		out += arg.String() + ", "
	}
	if len(self.Arguments) > 0 {
		out = out[:len(out)-2]
	}
	out += ")"
	return out
}

func NewFunctionCall(token *Token) *FunctionCall {
	return &FunctionCall{token.Str, []Expression{}, token}
}

func (self *FunctionCall) AppendArgument(arg Expression) {
	var args = append(self.Arguments, arg)
	self.Arguments = args
}

/*
TODO: Render function call

func (self FunctionCall) String() (out string) {
	out += self.Function
	out += "("
	for _, arg := range self.Arguments {

	}
	return self.Function + "(" + ")"
}
*/
