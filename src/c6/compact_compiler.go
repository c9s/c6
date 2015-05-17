package c6

import "c6/ast"

// import "fmt"

type CompactCompiler struct {
	Context *Context
}

func NewCompactCompiler(context *Context) *CompactCompiler {
	return &CompactCompiler{context}
}

func (compiler *CompactCompiler) CompileSimpleSelector(anySel ast.Selector) (out string) {
	return out
}

func (compiler *CompactCompiler) CompileCompoundSelector(compoundSelector *ast.CompoundSelector) (out string) {
	out = ""
	for _, sel := range *compoundSelector {
		out += compiler.CompileSimpleSelector(sel)
	}
	return out
}

func (compiler *CompactCompiler) CompileComplexSelector(sel *ast.ComplexSelector) (out string) {
	out = compiler.CompileCompoundSelector(sel.CompoundSelector)

	for _, item := range sel.ComplexSelectorItems {
		out += item.Combinator.String()
		out += compiler.CompileCompoundSelector(item.CompoundSelector)
	}
	return out
}

func (compiler *CompactCompiler) CompileComplexSelectorList(selectorList *ast.ComplexSelectorList) (out string) {
	out = ""
	for _, sel := range *selectorList {
		compiler.CompileComplexSelector(sel)
	}
	return out
}

func (compiler *CompactCompiler) CompileStatement(anyStm ast.Statement) string {

	switch stm := anyStm.(type) {
	case *ast.RuleSet:
		compiler.CompileComplexSelectorList(stm.Selectors)
	}

	return "stms"
}

func (compiler *CompactCompiler) Compile(any interface{}) (out string) {
	switch v := any.(type) {
	case []ast.Statement:
		for _, stm := range v {
			out += compiler.CompileStatement(stm)
		}
	}
	return out
}
