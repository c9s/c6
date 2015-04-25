package ast

type Property struct {
	Name PropertyName
	/**
	property value can be something like:

		`padding: 3px 3px;`
	*/
	Values []PropertyValue
}

/**
Property is one of the declaration
*/
func (self Property) IsDeclaration() {}

func (self Property) appendValue(value PropertyValue) {
	self.Values = append(self.Values, value)
}

type PropertyValue interface {
	CanBePropertyValue()
}

type PropertyName struct {
	String string
	// If there is an interpolation in the property name
	Interpolation bool
}
type ConstantValue string

func (self ConstantValue) CanBePropertyValue() {}

type Expression struct{}

func (self Expression) CanBePropertyValue() {}
