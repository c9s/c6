package ast

/*
DeclarationBlock.

.foo {
	property-name: property-value;
}
*/
type DeclarationBlock struct {
	Declarations []Declaration
}

/**
Append a Declaration
*/
func (self *DeclarationBlock) Append(decl Declaration) {
	self.Declarations = append(self.Declarations, decl)
}
