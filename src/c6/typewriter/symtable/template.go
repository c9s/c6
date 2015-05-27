package symtable

import "github.com/clipperhouse/typewriter"

// a convenience for passing values into templates; in MVC it'd be called a view model
type model struct {
	Type typewriter.Type
	Name string
	// these templates only ever happen to use one type parameter
	TypeParameter typewriter.Type
	Tag           typewriter.Tag
	// typewriter.TagValue
}

var templates = typewriter.TemplateSlice{symtableTemplate}

/*
type {{.Name}}SymTable map[string]*{{.Name}}

*/
var symtableTemplate = &typewriter.Template{
	Name: "SymTable",
	Text: `
func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

func (self *{{.Name}}) Set(name string, v *{{.TypeParameter.LongName}}) {
	self[name] = v
}

func (self {{.Name}}) Get(name string) (*{{.TypeParameter.LongName}}, bool) {
	if val, ok := self[name]; ok {
		return val, true
	}
	return nil, false
}

func (self *{{.Name}}) Merge(a *{{.Name}}) {
	for key, val := range *a {
		self[key] = val
	}
}

func (self *{{.Name}}) Has(name string) bool {
	if _, ok := self[name]; ok {
		return true
	} else {
		return false
	}
}
	
	`,
	TypeConstraint: typewriter.Constraint{Comparable: false},
}
