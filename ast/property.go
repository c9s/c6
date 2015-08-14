package ast

import "strings"

/**
An property may contains interpolation
*/
type Property struct {
	Name *PropertyName
	/**
	property value can be something like:

		`padding: 3px 3px;`
	*/
	Values []Expr
}

/**
Property is one of the declaration
*/
func (self Property) CanBeDeclaration() {}
func (self Property) CanBeStmt()        {}

func (self Property) AppendValue(value Expr) {
	self.Values = append(self.Values, value)
}

func (self Property) String() (out string) {
	out = self.Name.String() + ":"

	var items = []string{}

	for _, expr := range self.Values {
		items = append(items, expr.String())
	}
	out += strings.Join(items, " ")
	return out
}

type PropertyName struct {
	Name string
	// If there is an interpolation in the property name
	Interpolation bool
	Token         *Token
}

func (self PropertyName) String() string {
	return self.Name
}

func NewPropertyName(tok *Token) *PropertyName {
	return &PropertyName{tok.Str, tok.ContainsInterpolation, tok}
}

func NewProperty(nameTok *Token) *Property {
	return &Property{NewPropertyName(nameTok), []Expr{}}
}
