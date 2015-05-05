package ast

type Expression interface {
	// Evaluate() Value
	String() string
}

type UnaryExpression struct {
	Op   OpType
	Expr Expression
}

func NewUnaryExpression(op OpType, expr Expression) *UnaryExpression {
	return &UnaryExpression{op, expr}
}

func (self *UnaryExpression) Evaluate(symTable *SymTable) Value {
	var val Value = nil
	if expr, ok := self.Expr.(*BinaryExpression); ok {
		val = expr.Evaluate(symTable)
	} else if expr, ok := self.Expr.(*UnaryExpression); ok {
		val = expr.Evaluate(symTable)
	}

	// negative value
	if self.Op == OpSub {
		switch n := val.(type) {
		case *Number:
			n.Value = -n.Value
		case *Length:
			n.Value = -n.Value
		}
	}
	return val
}

func (self UnaryExpression) String() string {
	return "(" + self.Op.String() + " " + self.Expr.String() + ")"
}

type BinaryExpression struct {
	Op      OpType
	Left    Expression
	Right   Expression
	Grouped bool
}

func (self BinaryExpression) String() (out string) {
	out += "("
	out += self.Left.String()
	out += " " + self.Op.String() + " "
	out += self.Right.String()
	out += ")"
	return out
}

/**
The the divide expression will only be evaluated in the following 3 conditions:

	1. If the value, or any part of it, is stored in a variable or returned by a function.
	2. If the value is surrounded by parentheses.
	3. If the value is used as part of another arithmetic expression.

@see http://sass-lang.com/documentation/file.SASS_REFERENCE.html#division-and-slash
*/
func (self *BinaryExpression) IsCssSlash() bool {
	if self.Op == OpDiv {
		_, aok := self.Left.(*Length)
		_, bok := self.Right.(*Length)

		// it's not grouped, we should inflate it as string
		if aok && bok && self.Grouped == false {
			return true
		}
	}
	// otherwise we can divide the value
	return false
}

func (self *BinaryExpression) Evaluate(symTable *SymTable) Value {
	if self.IsCssSlash() {
		// return string object without quote
		return NewString(0, self.Left.(*Length).String()+"/"+self.Right.(*Length).String(), nil)
	}

	var lval Value = nil
	var rval Value = nil

	switch expr := self.Left.(type) {
	case *BinaryExpression:
		lval = expr.Evaluate(symTable)
	case *UnaryExpression:
		lval = expr.Evaluate(symTable)
	case *Number, *Length, *HexColor:
		lval = Value(expr)
	}
	switch expr := self.Right.(type) {
	case *UnaryExpression:
		rval = expr.Evaluate(symTable)
	case *BinaryExpression:
		rval = expr.Evaluate(symTable)
	case *Number, *Length, *HexColor:
		rval = Value(expr)
	}
	if lval != nil && rval != nil {
		return Compute(self.Op, lval, rval)
	}
	return nil
}

func NewBinaryExpression(op OpType, left Expression, right Expression, grouped bool) *BinaryExpression {
	return &BinaryExpression{op, left, right, grouped}
}
