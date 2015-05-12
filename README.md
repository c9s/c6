C6
===========================
Hate waiting for SASS compiling your stylesheets with Compass over 10 seconds
everytime?  C6 helps you write style sheets with efficiency.

C6 is a SASS compatible implementation written in Go. But wait! this is not only to
implement SASS, but also to improve the language for better consistency, syntax
and performance. And yes, this means we're free to accept new language
feature requests.

[![Build Status](https://travis-ci.org/c9s/c6.svg)](https://travis-ci.org/c9s/c6)
[![Coverage Status](https://coveralls.io/repos/c9s/c6/badge.svg)](https://coveralls.io/r/c9s/c6)
[![GoDoc](https://godoc.org/github.com/c9s/c6/src/c6?status.svg)](https://godoc.org/github.com/c9s/c6/src/c6)

[![wercker status](https://app.wercker.com/status/13aa03443c40dedeeabda923e1a95180/m "wercker status")](https://app.wercker.com/project/bykey/13aa03443c40dedeeabda923e1a95180)

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
  - [x] Nested properties.
  - [x] Comma-separated list
  - [x] Space-separated list
  - [x] `@if` , `@else` , `@else if`
  - [x] `@for`, `from`, `through` statement
  - [x] `@while`
  - [x] `@mixin`
  - [x] `@mixin` with arguments
  - [ ] `@include`
  - [x] Flags: `!default`, `!important`, `!optional`, `!global`
  - [x] Hex color
  - [x] Functions
  - [x] Vendor prefix properties
  - [x] MS filter.  `progid:DXImageTransform.Microsoft....`
  - [x] Variable names
  - [x] Variable assignment
  - [x] Time unit support. `s` second, `ms` ... etc
  - [x] Angle unit support.
  - [x] Resolution unit support.
  - [x] Unicode Range support: <https://developer.mozilla.org/en-US/docs/Web/CSS/unicode-range>
  - [x] Media Query
- [ ] Syntax
  - [ ] built-in `@import-once`
- [ ] Built-in Functions
  - .... to be listed
- [ ] Parser
  - [x] Parse `@import`
  - [x] Parse Expression
  - [x] Parse Space-Sep List
  - [x] Parse Comma-Sep List
  - [x] Parse Map (tests required)
  - [x] Parse Selector
  - [ ] Parse Selector with interpolation
  - [x] Parse RuleSet
  - [x] Parse DeclarationBlock
  - [x] Parse Variable Assignment Statement
  - [x] Parse PropertyName
  - [x] Parse PropertyName with interpolation
  - [x] Parse PropertyValue
  - [x] Parse PropertyValue with interpolation
  - [x] Parse conditions
  - [x] Parse `@media` statement
  - [ ] Parse Nested RuleSet
  - [x] Parse options: `!default`, `!global`, `!optional`
  - [ ] Parse CSS Hack for different browser (support more syntax sugar for this)
  - [ ] Parse `@font-face` block
  - [x] Parse `@if` statement
  - [x] Parse `@for` statement
  - [x] Parse `@while` statement
  - [ ] Parse `@mixin` statement
  - [ ] Parse `@include` statement
  - [ ] Parse `@function` statement
  - [ ] Parse `@return` statement
  - [ ] Parse keyword arguments for `@function`
  - [ ] Parse `@switch` statement
  - [ ] Parse `@case` statement
  - [ ] Parse `@use` statement

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
  - [x] Built-in color keyword table
  - [x] Hex Color computation
  - [x] Number operation: add, sub, mul, div
  - [x] Length operation: number operation for px, pt, em, rem, cm ...etc
  - [x] Expression evaluation
  - [x] Boolean expression evaluation
  - [x] Media Query conditions
  - [x] `@if` If Condition
  - [x] `@else if` If Else If
  - [x] `@else` else condition
  - [ ] `@while` statement node
  - [ ] `@each` statement node
  - [ ] `@function` statement node
  - [ ] `@mixin` statement node
  - [ ] `@include` statement node
  - [ ] `@return` statement node
- [ ] Runtime
  - [ ] HSL Color computation
- [ ] SASS Built-in Functions
  - [ ] RGB functions
    - [ ] `rgb($red, $green, $blue)`
    - [ ] `rgba($red, $green, $blue, $alpha)`
    - [ ] `red($color)`
    - [ ] `green($color)`
    - [ ] `blue($color)`
    - [ ] `mix($color1, $color2, [$weight])`
  - [ ] HSL Functions
    - [ ] `hsl($hue, $saturation, $lightness)`
    - [ ] `hsla($hue, $saturation, $lightness, $alpha)`
    - [ ] `hue($color)`
    - [ ] `saturation($color)`
    - [ ] `lightness($color)`
    - [ ] `adjust-hue($color, $degrees)`
    - [ ] `lighten($color, $amount)`
    - [ ] `darken($color, $amount)`
    - [ ] `saturate($color, $amount)`
    - [ ] `desaturate($color, $amount)`
    - [ ] `grayscale($color)`
    - [ ] `complement($color)`
    - [ ] `invert($color)`
  - [ ] Opacity Functions
    - [ ] `alpha($color) / opacity($color)`
    - [ ] `rgba($color, $alpha)`
    - [ ] `opacify($color, $amount) / fade-in($color, $amount)`
    - [ ] `transparentize($color, $amount) / fade-out($color, $amount)`
  - [ ] Other Color Functions
    - [ ] `adjust-color($color, [$red], [$green], [$blue], [$hue], [$saturation], [$lightness], [$alpha])`
    - [ ] `scale-color($color, [$red], [$green], [$blue], [$saturation], [$lightness], [$alpha])`
    - [ ] `change-color($color, [$red], [$green], [$blue], [$hue], [$saturation], [$lightness], [$alpha])`
    - [ ] `ie-hex-str($color)`
  - [ ] String Functions
    - [ ] `unquote($string)`
    - [ ] `quote($string)`
    - [ ] `str-length($string)`
    - [ ] `str-insert($string, $insert, $index)`
    - [ ] `str-index($string, $substring)`
    - [ ] `str-slice($string, $start-at, [$end-at])`
    - [ ] `to-upper-case($string)`
    - [ ] `to-lower-case($string)`
  - [ ] Number Functions
    - [ ] `percentage($number)`
    - [ ] `round($number)`
    - [ ] `ceil($number)`
    - [ ] `floor($number)`
    - [ ] `abs($number)`
    - [ ] `min($numbers…)`
    - [ ] `max($numbers…)`
    - [ ] `random([$limit])`
  - [ ] List Functions
    - [ ] `length($list)`
    - [ ] `nth($list, $n)`
    - [ ] `set-nth($list, $n, $value)`
    - [ ] `join($list1, $list2, [$separator])`
    - [ ] `append($list1, $val, [$separator])`
    - [ ] `zip($lists…)`
    - [ ] `index($list, $value)`
    - [ ] `list-separator(#list)`
  - [ ] Map Functions
    - [ ] `map-get($map, $key)`
    - [ ] `map-merge($map1, $map2)`
    - [ ] `map-remove($map, $keys…)`
    - [ ] `map-keys($map)`
    - [ ] `map-values($map)`
    - [ ] `map-has-key($map, $key)`
    - [ ] `keywords($args)`
  - [ ] Selector Functions
    - .... to be expanded ...


- [ ] CodeGen
  - [ ] NestedStyleCompiler
    - [ ] .... list ast nodes here ....

<!--
## Features

- [ ] import directory: https://github.com/sass/sass/issues/690
- [ ] import css as sass: https://github.com/sass/sass/issues/556
- [ ] import once: https://github.com/sass/sass/issues/139
- [ ] namespace and alias: https://github.com/sass/sass/issues/353
- [ ] `@use` directive: https://github.com/nex3/sass/issues/353#issuecomment-5146513 
- [ ] conditional import: https://github.com/sass/sass/issues/451
- [ ] `@sprite` syntax sugar
-->

## FAQ

### Why do you want to implement another SASS compiler?

The original SASS is written in Ruby and it's really slow, we take 11 seconds
to compile the stylesheets of our application, libsass is fast but it does not 
catch up the ruby sass implementation , and you can't compile them with Compass.
Since Go is simple, easy & fast, thus we implement SASS in Go to move faster.

Further, we want to design a new language that is compatible with SASS 3.4 or
even SASS 4.0 with more powerful features.

## Reference

SASS Reference <http://sass-lang.com/documentation/file.SASS_REFERENCE.html>


A feature check list from libsass:

- <https://github.com/sass/libsass/releases/tag/3.2.0>
- <https://github.com/sass/sass/issues/1094>

Grammar Ambiguity <https://www.facebook.com/cindylinz/posts/10202186527405801?hc_location=ufi>

Standards:

- CSS Syntax Module Level 3 <http://www.w3.org/TR/css-syntax-3>
- CSS Syntax Module Level 3 - Token Diagrams <http://www.w3.org/TR/css-syntax-3/#token-diagrams>

Frameworks, libraries:

- Bourbon <http://bourbon.io/>

- Marx <https://github.com/mblode/marx>

## Community

- Slack: we are on Slack channel `#c6`, go invite yourself here:
<https://docs.google.com/forms/d/11KpalZc6AUuQYf7vz531ys0pLEC9csjkC6QyC_hJQEg/viewform>

- Our Official Twitter Channel: <https://twitter.com/C6SASS>

## License

MPL License <https://www.mozilla.org/MPL/2.0/>

(MPL is like LGPL but with static/dynamic linking exception, which allows you
to either dynamic/static link this library without permissions)
