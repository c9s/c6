package symtable

import "github.com/clipperhouse/typewriter"

var templates = typewriter.TemplateSlice{set}

var set = &typewriter.Template{
	Name:           "SymTable",
	Text:           ``,
	TypeConstraint: typewriter.Constraint{Comparable: true},
}
