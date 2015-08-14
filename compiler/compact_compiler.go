package compiler

import "github.com/c9s/c6/ast"
import "github.com/c9s/c6/runtime"
import "bytes"
import "strings"

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
	Context      *runtime.Context
	ContextStack []runtime.Context
	Buffer       bytes.Buffer
	Compliant    int
	Indent       int
}

func NewCompactCompiler(context *runtime.Context) *CompactCompiler {
	return &CompactCompiler{
		Context:   context,
		Indent:    0,
		Compliant: 0,
	}
}

func Indent(level int) string {
	// two space
	return strings.Repeat("  ", level)
}

/*
 * Enable compliant
 */
func (compiler *CompactCompiler) EnableCompliant(compliant int) {
	compiler.Compliant |= compliant
}

/*
 * Disable compliant
 */
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

func (compiler *CompactCompiler) CompileDeclBlock(block *ast.DeclBlock) (out string) {
	out += "{"
	if block.Stmts != nil {
		for _, stm := range *block.Stmts {
			_ = stm
		}
	}
	out += "}"
	return out
}

func (compiler *CompactCompiler) CompileRuleSet(ruleset *ast.RuleSet) (out string) {
	compiler.Context.PushRuleSet(ruleset)
	out = compiler.CompileComplexSelectorList(ruleset.Selectors)
	out += " " + compiler.CompileDeclBlock(ruleset.Block)
	compiler.Context.PopRuleSet()
	return out
}

func (compiler *CompactCompiler) CompileStmt(anyStm ast.Stmt) string {

	switch stm := anyStm.(type) {
	case *ast.RuleSet:
		return compiler.CompileRuleSet(stm)
	case *ast.ImportStmt:
	case *ast.AssignStmt:
	}
	panic("Unsupported compilation")
}

func (compiler *CompactCompiler) CompileString(any interface{}) string {
	return compiler.Compile(any).String()
}

func (compiler *CompactCompiler) Compile(any interface{}) *bytes.Buffer {
	switch v := any.(type) {
	case ast.StmtList:
		for _, stm := range v {
			compiler.Buffer.WriteString(compiler.CompileStmt(stm))
		}
	case *ast.StmtList:
		for _, stm := range *v {
			compiler.Buffer.WriteString(compiler.CompileStmt(stm))
		}
	}
	return &compiler.Buffer
}
