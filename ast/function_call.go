package ast

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

func (self *FunctionCall) AppendArgument(arg Expr) {
	self.Arguments = append(self.Arguments, NewFunctionCallArgument(arg))
}

func NewFunctionCallWithToken(token *Token) *FunctionCall {
	return &FunctionCall{
		Ident:     token,
		Arguments: FunctionCallArguments{},
	}
}
