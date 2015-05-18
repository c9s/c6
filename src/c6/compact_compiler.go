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

func (compiler *CompactCompiler) DisableCompliant(compliant int) {
	compiler.Compliant = (compiler.Compliant | compliant) ^ compliant
}

func (compiler *CompactCompiler) HasCompliant(compliant int) bool {
	return (compiler.Compliant & compliant) > 0
}

func (compiler *CompactCompiler) CompileComplexSelector(sel *ast.ComplexSelector) (out string) {
	return sel.String()
}

func (compiler *CompactCompiler) CompileComplexSelectorList(selectorList *ast.ComplexSelectorList) string {
	return selectorList.String()
}

func (compiler *CompactCompiler) CompileDeclarationBlock(block *ast.DeclarationBlock) (out string) {
	out += "{"
	for _, stm := range block.Statements {
		_ = stm
	}
	out += "}"
	return out
}

func (compiler *CompactCompiler) CompileRuleSet(ruleset *ast.RuleSet) (out string) {
	out = compiler.CompileComplexSelectorList(ruleset.Selectors)
	out += " " + compiler.CompileDeclarationBlock(ruleset.Block)
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
