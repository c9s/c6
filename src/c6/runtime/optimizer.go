package runtime

import "c6/ast"

func OptimizeIfStmt(parentBlock *ast.Block, stm *ast.IfStmt) {

	// TODO: select passed condition and merge block
	// try to simplify the condition without context and symbol table

	var mergeBlock = false
	var ignoreBlock = false
	var val = EvaluateExprInBooleanContext(stm.Condition, nil)
	// check if the expression is evaluated
	if IsValue(val) {
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
