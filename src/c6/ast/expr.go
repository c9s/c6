package ast

type Expression interface {
	String() string
}

type UnaryExpression struct {
	Op   *Op
	Expr Expression
}

func NewUnaryExpression(op *Op, expr Expression) *UnaryExpression {
	return &UnaryExpression{op, expr}
}

func (self UnaryExpression) String() string {
	return "(" + self.Op.String() + " " + self.Expr.String() + ")"
}

type BinaryExpression struct {
	Op      *Op
	Left    Expression
	Right   Expression
	Grouped bool
}

func (self BinaryExpression) String() (out string) {
	out += "("
	out += self.Left.String()
	if self.Op != nil {
		out += " " + self.Op.String() + " "
	}
	out += self.Right.String()
	out += ")"
	return out
}

/*
If any of left or right is variable, than it's constant expression. this is
used to eliminate simple expression when the parser is parsing...

Please note thist method does not test CSS slash, the caller should handle by itself.

This works for both boolean evaluation and arithmetic evaluation.
*/
func (self BinaryExpression) IsConstantExpression() bool {
	_, ok1 := self.Left.(*Variable)
	_, ok2 := self.Right.(*Variable)
	return ok1 || ok2
}

/**
The the divide expression will only be evaluated in the following 3 conditions:

	1. If the value, or any part of it, is stored in a variable or returned by a function.
	2. If the value is surrounded by parentheses.
	3. If the value is used as part of another arithmetic expression.

This method needs to be called on the top caller to prevent unexpected result.

@see http://sass-lang.com/documentation/file.SASS_REFERENCE.html#division-and-slash
*/
func (self *BinaryExpression) IsCssSlash() bool {
	if self.Op.Type == T_DIV {
		_, aok := self.Left.(*Number)
		_, bok := self.Right.(*Number)

		// it's not grouped, we should inflate it as string
		if aok && bok && self.Grouped == false {
			return true
		}
	}
	// otherwise we can divide the value
	return false
}

func NewBinaryExpression(op *Op, left Expression, right Expression, grouped bool) *BinaryExpression {
	return &BinaryExpression{op, left, right, grouped}
}
