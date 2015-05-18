package c6

import "c6/ast"

// import "fmt"

const (
	_             = iota
	CSS3Compliant = 1 << (2 * iota)
	CSS4Compliant
	IE7Compliant
	IE8Compliant
	IE9Compliant
	IE10Compliant
	WebKitCompliant
	MozCompliant // optimize for -moz
)

type CompactCompiler struct {
	Context   *Context
	Compliant int
}

func NewCompactCompiler(context *Context) *CompactCompiler {
	return &CompactCompiler{
		Context: context,
	}
}

func (compiler *CompactCompiler) EnableCompliant(compliant int) {
	compiler.Compliant |= compliant
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

func (compiler *CompactCompiler) CompileDeclarationBlock(block *ast.DeclarationBlock) (out string) {
	return out
}

func (compiler *CompactCompiler) CompileRuleSet(ruleset *ast.RuleSet) (out string) {
	out = compiler.CompileComplexSelectorList(ruleset.Selectors)
	out += compiler.CompileDeclarationBlock(ruleset.Block)
	return out
}

func (compiler *CompactCompiler) CompileStatement(anyStm ast.Statement) string {

	switch stm := anyStm.(type) {
	case *ast.RuleSet:
		return compiler.CompileRuleSet(stm)
	case *ast.ImportStatement:
	case *ast.VariableAssignment:
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
