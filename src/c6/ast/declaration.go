package ast

/*
A declaration can be a property or a ruleset
*/
type Declaration interface {
	IsDeclaration()
}
