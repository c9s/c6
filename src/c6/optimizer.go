package c6

import "github.com/c9s/c6/src/c6/ast"

func OptimizeIfStatement(parentBlock *ast.Block, stm *ast.IfStatement) {

	// TODO: select passed condition and merge block
	// try to simplify the condition without context and symbol table

	var mergeBlock = false
	var ignoreBlock = false
	var val = EvaluateExpressionInBooleanContext(stm.Condition, nil)
	// check if the expression is evaluated
	if IsConstantValue(val) {
		if b, ok := val.(ast.BooleanValue); ok {
			if b.Boolean() {
				mergeBlock = true
			} else {
				ignoreBlock = true
			}
		}
	}

	if mergeBlock {
		// TODO: merge with subblock with the current block
	} else if ignoreBlock {

	}

}
