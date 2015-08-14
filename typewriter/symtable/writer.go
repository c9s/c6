package symtable

import (
	"github.com/clipperhouse/typewriter"
	"io"
	"log"
)

func init() {
	err := typewriter.Register(NewSymTableWriter())
	if err != nil {
		panic(err)
	}
}

type SymTableWriter struct{}

func NewSymTableWriter() *SymTableWriter {
	return &SymTableWriter{}
}

func (sw *SymTableWriter) Name() string {
	return "symtable"
}

func (sw *SymTableWriter) Imports(t typewriter.Type) []typewriter.ImportSpec {
	imports := []typewriter.ImportSpec{
		{Path: "github.com/c9s/c6/ast"},
	}
	return imports
}

func (sw *SymTableWriter) Write(w io.Writer, t typewriter.Type) error {
	tag, found := t.FindTag(sw)

	if !found {
		// nothing to be done
		return nil
	}

	log.Printf("Found tag %+v\n", tag)

	license := `// SymTableWriter is a modification of https://github.com/deckarep/golang-set
// The MIT License (MIT)
// Copyright (c) 2015 Yo-An Lin (yoanlin93@gmail.com)
`
	if _, err := w.Write([]byte(license)); err != nil {
		return err
	}

	log.Printf("Finding template by tag %+v\n", tag)
	tmpl, err := templates.ByTag(t, tag)
	if err != nil {
		return err
	}

	m := model{
		Type:          t,
		Name:          t.Name,
		TypeParameter: tag.Values[0].TypeParameters[0],
		Tag:           tag,
	}

	log.Println("Rendering templates")
	if err := tmpl.Execute(w, m); err != nil {
		return err
	}

	log.Println("Done!")
	return nil
}
