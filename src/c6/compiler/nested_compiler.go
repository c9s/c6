package compiler

import "github.com/c9s/c6/src/c6/ast"

type Compiler interface {
	CompileBlock(block *ast.Block) string
}

type NestedStyleCompiler struct {
	Indent int
	Output string
}

func NewNestedStyleCompiler() *NestedStyleCompiler {
	return &NestedStyleCompiler{}
}

func (self *NestedStyleCompiler) CompileProperty(property *ast.Property) {

}

func (self *NestedStyleCompiler) CompileSeletors(selectors *ast.Selector) {

}

func (self *NestedStyleCompiler) CompileRuleSet(ruleset *ast.RuleSet) {

}

func (self *NestedStyleCompiler) CompileBlock(block *ast.Block) {

}
