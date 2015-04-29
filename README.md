C6
===========================

C6 is a SASS implementation in Go.

This is not just to implement SASS, but also to improve the language for better
consistency, syntax and performance. And yes, this means we're free to accept any new
language feature requests.


## Setup

Setup GOPATH:

    source goenv

Setup dependencies:

    source goinstall

Run tests:

    go test c6
    go test -v c6
    go test -v c6/ast

To run specific test

    go test -run TestParser -x -v c6

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


## Features

- [ ] import directory: https://github.com/sass/sass/issues/690
- [ ] import css as sass: https://github.com/sass/sass/issues/556
- [ ] import once: https://github.com/sass/sass/issues/139
- [ ] namespace and alias: https://github.com/sass/sass/issues/353
- [ ] `@use` directive: https://github.com/nex3/sass/issues/353#issuecomment-5146513 
- [ ] conditional import: https://github.com/sass/sass/issues/451



## Reference

A feature check list from libsass:

- https://github.com/sass/libsass/releases/tag/3.2.0
- https://github.com/sass/sass/issues/1094


## License

MPL License <https://www.mozilla.org/MPL/2.0/>
