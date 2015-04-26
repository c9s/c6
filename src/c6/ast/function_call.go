package ast

type FunctionCall struct {
	Function  string
	Arguments []Expression
	Token     *Token
}

func NewFunctionCall(token *Token) *FunctionCall {
	return &FunctionCall{token.Str, []Expression{}, token}
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
