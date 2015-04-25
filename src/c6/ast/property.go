package ast

type Property struct {
	Name PropertyName
	/**
	property value can be something like:

		`padding: 3px 3px;`
	*/
	Values []Expression
}

/**
Property is one of the declaration
*/
func (self Property) IsDeclaration() {}

func (self Property) AppendValue(value Expression) {
	self.Values = append(self.Values, value)
}

type PropertyName struct {
	String string
	// If there is an interpolation in the property name
	Interpolation bool
	Token         Token
}

type Expression interface {
	CanBeExpression()
}

type UnaryExpression struct {
	Value interface{}
	Token Token
}

func (self UnaryExpression) CanBeExpression() {}

type BinaryExpression struct {
	Left  Expression
	Right Expression
	Op    string
}

func (self BinaryExpression) CanBeExpression() {}

type ConstantString struct {
	Constant string
	Token    Token
}

func (self ConstantString) CanBeExpression() {}
