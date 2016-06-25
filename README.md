C6
===========================
Hate waiting for Sass compiling your stylesheets with Compass over 10 seconds
everytime?  C6 helps you write style sheets with efficiency.

[![Build Status](https://travis-ci.org/c9s/c6.svg)](https://travis-ci.org/c9s/c6)
[![Coverage Status](https://coveralls.io/repos/c9s/c6/badge.svg)](https://coveralls.io/r/c9s/c6)
[![wercker status](https://app.wercker.com/status/13aa03443c40dedeeabda923e1a95180/m "wercker status")](https://app.wercker.com/project/bykey/13aa03443c40dedeeabda923e1a95180)

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
  - [x] `@include`
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
- [ ] Built-in Functions
  - .... to be listed
- [ ] Parser
  - [x] Parse `@import`
  - [x] Parse Expr
  - [x] Parse Space-Sep List
  - [x] Parse Comma-Sep List
  - [x] Parse Map (tests required)
  - [x] Parse Selector
  - [ ] Parse Selector with interpolation
  - [x] Parse RuleSet
  - [x] Parse DeclBlock
  - [x] Parse Variable Assignment Stmt
  - [x] Parse PropertyName
  - [x] Parse PropertyName with interpolation
  - [x] Parse PropertyValue
  - [x] Parse PropertyValue with interpolation
  - [x] Parse conditions
  - [x] Parse `@media` statement
  - [x] Parse Nested RuleSet
  - [x] Parse Nested Properties
  - [x] Parse options: `!default`, `!global`, `!optional`
  - [ ] Parse CSS Hack for different browser (support more syntax sugar for this)
  - [x] Parse `@font-face` block
  - [x] Parse `@if` statement
  - [x] Parse `@for` statement
  - [x] Parse `@while` statement
  - [x] Parse `@mixin` statement
  - [x] Parse `@include` statement
  - [x] Parse `@function` statement
  - [x] Parse `@return` statement
  - [x] Parse `@extend` statement
  - [x] Parse keyword arguments for `@function`
  - [ ] Parse `@switch` statement
  - [ ] Parse `@case` statement
  - [ ] Parse `@use` statement
- [ ] Building AST
  - [x] RuleSet
  - [x] DeclBlock
  - [x] PropertyName
  - [x] PropertyValue
  - [x] Comma-Separated List
  - [x] Space-Separated List
  - [x] Basic Exprs
  - [x] FunctionCall
  - [x] Expr with interpolation
  - [x] Variable statements
  - [x] Built-in color keyword table
  - [x] Hex Color computation
  - [x] Number operation: add, sub, mul, div
  - [x] Length operation: number operation for px, pt, em, rem, cm ...etc
  - [x] Expr evaluation
  - [x] Boolean expression evaluation
  - [x] Media Query conditions
  - [x] `@if` If Condition
  - [x] `@else if` If Else If
  - [x] `@else` else condition
  - [x] `@while` statement node
  - [x] `@function` statement node
  - [x] `@mixin` statement node
  - [x] `@include` statement node
  - [x] `@return` statement node
  - [ ] `@each` statement node

- [ ] Runtime
  - [ ] HSL Color computation
  - [ ] Function Call Invoke mech
  - [ ] Mixin Include
  - [ ] Import

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
  - [ ] CompactCompiler
    - [ ] CompileCssImportStmt: `@import url(...);`
    - [ ] CompileRuleSet
    - [ ] CompileSelectors
      - [ ] CoimpileParentSelector
    - [ ] CompileSubRuleSet
    - [ ] CompileCommentBlock
    - [ ] CompileDeclBlock
    - [ ] CompileMediaQuery: `@media`
    - [ ] CompileSupportrQuery: `@support`
    - [ ] CompileFontFace: `@support`
    - [ ] CompileForStmt
    - [ ] CompileIfStmt
      - [ ] CompileElseIfStmt
    - [ ] CompileWhileStmt
    - [ ] CompileEachStmt
    - [ ] ... list more ast nodes here ...

- [ ] Syntax
  - [ ] built-in `@import-once`

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


## Self Benchmarking

[![BenchViz](https://raw.githubusercontent.com/c9s/c6/master/benchmarks/summary.svg)](https://raw.githubusercontent.com/c9s/c6/master/benchmarks/summary.svg)

## FAQ

### Why do you want to implement another Sass compiler?

The original Sass is written in Ruby and it's really slow, we take 11 seconds
to compile the stylesheets of our application, libsass is fast but it does not 
catch up the ruby sass implementation , and you can't compile them with Compass.
Since Go is simple, easy & fast, thus we implement Sass in Go to move faster.

Further, we want to design a new language (ECSS temporarily named) that is compatible with Sass 3.4 or
even Sass 4.0 with more powerful features.

## Reference

Sass Reference <http://sass-lang.com/documentation/file.SASS_REFERENCE.html>


A feature check list from libsass:

- <https://github.com/sass/libsass/releases/tag/3.2.0>
- <https://github.com/sass/sass/issues/1094>
- <https://github.com/sass/sass/issues/739> Dynamic dependencies

Grammar Ambiguity:

- <https://www.facebook.com/cindylinz/posts/10202186527405801?hc_location=ufi>
- <https://www.facebook.com/yoan.lin/posts/10152968537931715?_rdr>

Standards:

- CSS Syntax Module Level 3 <http://www.w3.org/TR/css-syntax-3>
- CSS 3 Selector <http://www.w3.org/TR/css3-selectors/#grouping>
- CSS Font <http://www.w3.org/TR/css3-fonts/#basic-font-props>
- Selectors API <http://www.w3.org/TR/selectors-api/>
- At-Page Rule <http://dev.w3.org/csswg/css-page-3/#at-page-rule>
- Railroad diagram <https://github.com/tabatkins/railroad-diagrams>
- CSS 2.1 Grammar <http://www.w3.org/TR/CSS21/grammar.html>

Sass/CSS Frameworks, libraries:

- Bourbon <http://bourbon.io/>
- Marx <https://github.com/mblode/marx>
- FormHack <http://formhack.io/>
- Susy <http://susy.oddbird.net/>
- Gumby <http://www.gumbyframework.com/>

Articles:

- Logic in media queries - <https://css-tricks.com/logic-in-media-queries/>

Go related:

- <https://docs.google.com/document/d/1y4mYe8Sk9jCPze6AyygxrC7j1sqEgwAycezpaE6JnC8/edit#heading=h.xgr9eiceu2ca>


## Contribution

Please open up an issue on GitHub before you put a lot efforts on pull request.

The code submitting to PR must be filtered with `gofmt`

## Community

- Slack: we are on Slack channel `#c6`, go invite yourself here:
<https://docs.google.com/forms/d/11KpalZc6AUuQYf7vz531ys0pLEC9csjkC6QyC_hJQEg/viewform>

- Our Official Twitter Channel: <https://twitter.com/C6SASS>

## License

MPL License <https://www.mozilla.org/MPL/2.0/>

(MPL is like LGPL but with static/dynamic linking exception, which allows you
to either dynamic/static link this library without permissions)


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/c9s/c6/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

