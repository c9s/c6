package symtable

import (
	"github.com/clipperhouse/typewriter"
	"io"
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
	return "set"
}

func (sw *SymTableWriter) Imports(t typewriter.Type) (result []typewriter.ImportSpec) {
	// none
	return result
}

func (sw *SymTableWriter) Write(w io.Writer, t typewriter.Type) error {
	tag, found := t.FindTag(sw)

	if !found {
		// nothing to be done
		return nil
	}

	license := `// SymTableWriter is a modification of https://github.com/deckarep/golang-set
// The MIT License (MIT)
// Copyright (c) 2015 Yo-An Lin (yoanlin93@gmail.com)
`
	if _, err := w.Write([]byte(license)); err != nil {
		return err
	}

	tmpl, err := templates.ByTag(t, tag)

	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, t); err != nil {
		return err
	}
	return nil
}
