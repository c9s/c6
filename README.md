C6
===========================
Hate waiting for SASS compiling your stylesheets with Compass over 10 seconds
everytime?  C6 helps you write style sheets with efficiency.

C6 is a SASS implementation written in Go. But wait! this is not only to
implement SASS, but also to improve the language for better consistency, syntax
and performance. And yes, this means we're free to accept new language
feature requests.

[![Build Status](https://travis-ci.org/c9s/c6.svg)](https://travis-ci.org/c9s/c6)
[![Coverage Status](https://coveralls.io/repos/c9s/c6/badge.svg)](https://coveralls.io/r/c9s/c6)
[![GoDoc](https://godoc.org/github.com/c9s/c6/src/c6?status.svg)](https://godoc.org/github.com/c9s/c6/src/c6)

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
  - [x] Ruleset
  - [x] Sub-ruleset
  - [x] Interpolation
  - [x] Property name
  - [x] Property value list
  - [x] Comma-separated list
  - [x] Space-separated list
  - [x] Hex color
  - [x] Functions
  - [x] Vendor prefix properties
  - [x] MS filter.  `progid:DXImageTransform.Microsoft....`
  - [x] Variable names
  - [x] Variable assignment
  - [x] `s` second, `ms` unit lexing support.
  - [ ] Media Query
- [ ] Syntax
  - [ ] built-in `@import-once`
- [ ] Built-in Functions
  - .... to be listed
- [ ] Parser
  - [x] Parse `@import`
  - [x] Parse Expression
  - [x] Parse Space-Sep List
  - [x] Parse Comma-Sep List
  - [-] Parse Map (tests required)
  - [ ] Parse Condition
  - [x] Parse Selector
  - [ ] Parse Selector with interpolation
  - [x] Parse RuleSet
  - [x] Parse DeclarationBlock
  - [x] Parse Variable Assignment Statement
  - [x] Parse PropertyName
  - [ ] Parse PropertyName with interpolation
  - [-] Parse PropertyValue
  - [-] Parse PropertyValue with interpolation
  - [ ] Parse Nested RuleSet
  - [ ] Parse options: `!default`, `!global`, `!optional`
  - [ ] Parse CSS Hack for different browser (support more syntax sugar for this)
  - [ ] Parse `@if`
  - [ ] Parse `@mixin`
  - [ ] Parse `@include`
  - [ ] Parse `@function`
  - [ ] Parse `@media`
- [ ] Building AST
  - [x] RuleSet
  - [x] DeclarationBlock
  - [x] PropertyName
  - [x] PropertyValue
  - [x] Comma-Separated List
  - [x] Space-Separated List
  - [x] Basic Expressions
  - [x] FunctionCall
  - [x] Expression with interpolation
  - [x] Variable statements
  - [ ] If Condition
  - [ ] If Else If, Else Condition
  - [ ] Built-in color keyword table
  - [ ] Hex Color computation
  - [ ] HSL Color computation
  - [ ] Number operation: add, sub, mul, div
  - [ ] Length operation: number operation for px, pt, em, rem, cm ...etc
  - [ ] Expression evaluation
  - [ ] Media Query conditions
- [ ] CodeGen
  - [ ] NestedStyleCompiler
    - [ ] .... list ast nodes here ....

## Features

- [ ] import directory: https://github.com/sass/sass/issues/690
- [ ] import css as sass: https://github.com/sass/sass/issues/556
- [ ] import once: https://github.com/sass/sass/issues/139
- [ ] namespace and alias: https://github.com/sass/sass/issues/353
- [ ] `@use` directive: https://github.com/nex3/sass/issues/353#issuecomment-5146513 
- [ ] conditional import: https://github.com/sass/sass/issues/451
- [ ] `@sprite` syntax sugar



## Reference

A feature check list from libsass:

- https://github.com/sass/libsass/releases/tag/3.2.0
- https://github.com/sass/sass/issues/1094


## License

MPL License <https://www.mozilla.org/MPL/2.0/>

(MPL is like LGPL but with static/dynamic linking exception, which allows you
to either dynamic/static link this library without permissions)
