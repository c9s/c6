package ast

import (
	"bytes"
	"fmt"
)

type Scope struct {
	Parent  *Scope
	Objects map[string]*Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Objects: make(map[string]*Object, 4)}
}

func (s *Scope) Lookup(name string) *Object {
	return s.Objects[name]
}

func (s *Scope) Insert(obj *Object) (alt *Object) {
	if alt = s.Objects[obj.Name]; alt == nil {
		s.Objects[obj.Name] = obj
	}
	return
}

// Debugging support
func (s *Scope) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Scope %p {", s)
	if s != nil && len(s.Objects) > 0 {
		fmt.Fprintln(&buf)
		for _, obj := range s.Objects {
			fmt.Fprintf(&buf, "\t%s %s\n", obj.Kind, obj.Name)
		}
	}
	fmt.Fprintf(&buf, "}\n")
	return buf.String()
}

type Object struct {
	Name string
	Kind ObjectKind
	Decl interface{}
}

// ObjectKind describes what an object represents.
type ObjectKind int

// The list of possible Object kinds.
const (
	Bad   ObjectKind = iota // for error handling
	Pkg                     // package
	Var                     // variable
	Fun                     // function
	Mixin                   // mixin
)

var objectKindStrings = [...]string{
	Bad:   "bad",
	Pkg:   "package",
	Var:   "var",
	Fun:   "func",
	Mixin: "mixin",
}

func (kind ObjectKind) String() string { return objectKindStrings[kind] }

// This function creates a new object of a given kind and name.
func NewObject(kind ObjectKind, name string) *Object {
	return &Object{Kind: kind, Name: name}
}
