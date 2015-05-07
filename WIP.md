======

VariableAssignment Flag support:

- [ ] Add Flag struct to VariableAssignment struct
- [ ] Parse Flag keywords and push to VariableAssignment struct

CharsetStatement support

- [ ] Add CharsetStatement struct
- [ ] Add Encoding field
- [ ] Add ParseCharsetStatement method to parser.

SymbolTable

- [ ] Register parsed variable to the scope symbol table.
  - [ ] RuleSet symbol table
  - [ ] Global symbol table
- [ ] Add symbol table lookup method to the expression evaluator.
  - Add type switch case for ast.Variable struct

Optimizer

- [ ] Constant Value elimination optimizer for VariableAssignment.
- [ ] Call IfStatementOptimizer after the if statement is parsed.

CSS Slash and Divide

- [ ] Review Declaration String() interface.
- [ ] Test simple ruleset output.
- [ ] Add test utility function that accept: {input scss} and {output css}.
- [ ] Add expr stringer test case for `font: 12px/20px`.
- [ ] Add expr stringer test case for expressions like 12px/20px + 13px.

Nested properties

- [ ] Allow declaration block after the colon of property name.
- [ ] Allow declaration block after the property value.
- [ ] `lexPropertyValue` should check if there is another '{' token, then we should go to `lexStatement` state.

