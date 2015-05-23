package c6

/**
We will cache the compiled ast.Block in the map,
*/
type ASTFileMap map[string]interface{}

var ASTCache ASTFileMap = ASTFileMap{}
