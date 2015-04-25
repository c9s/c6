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

func (self ConstantString) CanBeExpression() {}
