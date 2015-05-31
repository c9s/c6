======

Parser Information

- [ ] File Level Scope should record funtion decl, mixin decl, var decl
- [ ] Create new scope object when entering mixin, function or a declaration
  block of a ruleset.
- [ ] When parsing expression, we should lookup the related symbol and connect
  the parsed factors.  For example, connect `$a` to the parent scope `$a`
  declaration. this is required for compiler to evaluate the result

- [ ] The ast.Scope contains the declaration, not the evaludated value.
- [ ] compiler.Scope is somewhere that might contains the evaluated value.

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

- [x] Constant Value elimination optimizer for AssignStmt.
- [ ] Call IfStmtOptimizer after the if statement is parsed.

CSS Slash and Divide

- [x] Review Declaration String() interface.
- [ ] Test simple ruleset output.
- [ ] Add test utility function that accept: {input scss} and {output css}.
- [ ] Add expr stringer test case for `font: 12px/20px`.
- [ ] Add expr stringer test case for expressions like 12px/20px + 13px.

Nested properties

- [ ] Allow declaration block after the colon of property name.
- [ ] Allow declaration block after the property value.
- [ ] `lexPropertyValue` should check if there is another '{' token, then we should go to `lexStmt` state.

`@my` statement

- [ ] Declare variable in the specific scope.


