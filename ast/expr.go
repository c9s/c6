package ast

type Expr interface {
	String() string
}

type UnaryExpr struct {
	Op   *Op
	Expr Expr
}

func NewUnaryExpr(op *Op, expr Expr) *UnaryExpr {
	return &UnaryExpr{op, expr}
}

func (self UnaryExpr) String() string {
	if self.Op != nil {
		return self.Op.String() + self.Expr.String()
	}
	return self.Expr.String()
}

type BinaryExpr struct {
	Op      *Op
	Left    Expr
	Right   Expr
	Grouped bool
}

func (self BinaryExpr) String() string {
	if self.Op == nil {
		panic("Missing operator")
	}

	var out = self.Left.String() + self.Op.String() + self.Right.String()
	if self.Grouped {
		return "(" + out + ")"
	}
	return out
}

/*
If any of left or right is variable, than it's constant expression. this is
used to eliminate simple expression when the parser is parsing...

Please note thist method does not test CSS slash, the caller should handle by itself.

This works for both boolean evaluation and arithmetic evaluation.
*/
func (self BinaryExpr) IsSimpleExpr() bool {
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
func (self *BinaryExpr) IsCssSlash() bool {
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

func NewBinaryExpr(op *Op, left Expr, right Expr, grouped bool) *BinaryExpr {
	return &BinaryExpr{op, left, right, grouped}
}
