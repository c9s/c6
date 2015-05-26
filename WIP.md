======

SymbolTable

- [x] Register parsed variable to the scope symbol table.
  - [x] RuleSet symbol table
  - [x] Global symbol table
- [x] Add symbol table lookup method to the expression evaluator.
  - [x] Add type switch case for ast.Variable struct
- [ ] Remove symbol table from DeclarationBlock since we need new context to evaluate the blocks.
- [ ] The context needs to be a stack, so we can push new context into the stack to evaluate the values.
- [ ] Only constants can be registered to the DeclarationBlock

Function Evaluation

Evaluating function calls, we:

1. Fetch the funtion definition 
2. Push the arguments into the stack.
3. Create a new context for the function call
4. Create a symbol table for the function call.
5. Pop the arguments from the stack into a new content.
6. Push the return value onto the stack.
7. Go back to the caller.
8. Pop the returning result from the stack.






Optimizer

- [x] Constant Value elimination optimizer for VariableAssignment.
- [ ] Call IfStatementOptimizer after the if statement is parsed.

CSS Slash and Divide

- [x] Review Declaration String() interface.
- [ ] Test simple ruleset output.
- [ ] Add test utility function that accept: {input scss} and {output css}.
- [ ] Add expr stringer test case for `font: 12px/20px`.
- [ ] Add expr stringer test case for expressions like 12px/20px + 13px.

Nested properties

- [ ] Allow declaration block after the colon of property name.
- [ ] Allow declaration block after the property value.
- [ ] `lexPropertyValue` should check if there is another '{' token, then we should go to `lexStatement` state.

`@my` statement

- [ ] Declare variable in the specific scope.


