package ast

import "fmt"

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
	Op    OpType
	Left  Expression
	Right Expression
}

func (self BinaryExpression) String() (out string) {
	out += "("
	out += self.Left.String()
	out += " " + self.Op.String() + " "
	out += self.Right.String()
	out += ")"
	return out
}

func (self *BinaryExpression) Evaluate(symTable *SymTable) Value {
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

func NewBinaryExpression(op OpType, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{op, left, right}
}

func Compute(op OpType, a Value, b Value) Value {
	switch op {
	case OpAdd:
		switch ta := a.(type) {
		case *Number:
			switch tb := b.(type) {
			case *Number:
				return NumberAddNumber(ta, tb)
			case *HexColor:
				return HexColorAddNumber(tb, ta)
			}
		case *Length:
			switch tb := b.(type) {
			case *Length:
				return LengthAddLength(ta, tb)
			}
		case *HexColor:
			switch tb := b.(type) {
			case *Number:
				return HexColorAddNumber(ta, tb)
			}
		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorAddNumber(ta, tb)
			}
		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorAddNumber(ta, tb)
			}
		}
	case OpSub:
		switch ta := a.(type) {

		case *Number:
			switch tb := b.(type) {
			case *Number:
				return NumberSubNumber(ta, tb)
			}

		case *Length:
			switch tb := b.(type) {
			case *Length:
				val := LengthSubLength(ta, tb)
				fmt.Printf("Substracted value: %+v\n", val)
				return val
			}

		case *HexColor:
			switch tb := b.(type) {
			case *Number:
				return HexColorSubNumber(ta, tb)
			}

		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorSubNumber(ta, tb)
			}

		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorSubNumber(ta, tb)
			}
		}
	case OpMul:
		switch ta := a.(type) {

		case *Length:
			switch tb := b.(type) {
			case *Length:
				return LengthMulLength(ta, tb)
			case *Number:
				return LengthMulNumber(ta, tb)
			}

		case *RGBColor:
			switch tb := b.(type) {
			case *Number:
				return RGBColorMulNumber(ta, tb)
			}

		case *RGBAColor:
			switch tb := b.(type) {
			case *Number:
				return RGBAColorMulNumber(ta, tb)
			}
		}
	}
	return nil
}
