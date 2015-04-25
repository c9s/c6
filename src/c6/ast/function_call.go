package ast

type FunctionCall struct {
	Function  string
	Token     Token
	Arguments []Expression
}

func NewFunctionCall(name string, token Token) *FunctionCall {
	return &FunctionCall{name, token, []Expression{}}
}
func (self FunctionCall) CanBeExpression() {}

func (self FunctionCall) AppendArgument(arg Expression) {
	self.Arguments = append(self.Arguments, arg)
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
