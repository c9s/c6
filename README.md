C6
===========================

A SASS implementation in Go.

This is not just to implement SASS, but also to improve the language for better
consistency, syntax and performance. And yes, this means we're free to accept any new
language feature requests.

## Working in progress


- [ ] Lexing
  - [x] `@import`
  - [x] simple selector.
    - [x] type selector.
    - [x] child selector.
    - [x] attribute selector.
    - [x] adjacent selector.
    - [x] descendant selector.
    - [x] class selector.
    - [x] ID selector.
  - [x] ruleset
  - [x] sub-ruleset
  - [x] interpolation
  - [x] property name
  - [x] property value list
  - [x] comma-separated list
  - [x] space-separated list
  - [x] hex color
  - [x] functions
  - [x] vendor prefix properties
  - [x] MS filter.  `progid:DXImageTransform.Microsoft....`
  - [x] variable names
  - [x] variable assignment
- [ ] Syntax
  - [x] `@import`
  - [ ] `@if`
  - [ ] `@mixin`
  - [ ] `@include`
  - [ ] `@function`
  - [ ] `@media`
  - [ ] built-in `@import-once`
- [ ] Built-in Functions
  - .... to be listed
- [ ] Building AST
  - [x] RuleSet
  - [x] DeclarationBlock
  - [x] PropertyName
  - [x] PropertyValue
  - [x] Comma-Separated List
  - [x] Space-Separated List
  - [x] Basic Expressions
  - [ ] Variable statements
- [ ] NestedStyleCompiler

