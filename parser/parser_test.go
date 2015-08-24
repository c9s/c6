package parser

import (
	"io/ioutil"
	"testing"

	"github.com/c9s/c6/ast"
	"github.com/c9s/c6/runtime"

	"github.com/stretchr/testify/assert"
)

func RunParserTest(code string) *ast.StmtList {
	var p = NewParser(runtime.NewContext())
	return p.ParseScss(code)
}

func TestParserGetFileType(t *testing.T) {
	matrix := map[uint]string{
		UnknownFileType: ".css",
		ScssFileType:    ".scss",
		SassFileType:    ".sass",
		EcssFileType:    ".ecss",
	}

	for k, v := range matrix {
		assert.Equal(t, k, getFileTypeByExtension(v))
	}

}

func TestParserParseFile(t *testing.T) {
	testPath := "test/file.scss"
	bs, _ := ioutil.ReadFile(testPath)
	p := NewParser(runtime.NewContext())
	err := p.ParseFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if e := string(bs); e != p.Content {
		t.Fatalf("got: %s wanted: %s", p.Content, e)
	}

	if e := testPath; e != p.File.FileName {
		t.Fatalf("got: %s wanted: %s", p.File.FileName, e)
	}
}

func TestParserEmptyRuleSetWithUniversalSelector(t *testing.T) {
	var stmts = RunParserTest(`* { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithClassSelector(t *testing.T) {
	var stmts = RunParserTest(`.first-name { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithIdSelector(t *testing.T) {
	var stmts = RunParserTest(`#myId { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithTypeSelector(t *testing.T) {
	var stmts = RunParserTest(`div { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorAttributeNameOnly(t *testing.T) {
	var stmts = RunParserTest(`[href] { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorPrefixMatch(t *testing.T) {
	var stmts = RunParserTest(`[href^=http] { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserEmptyRuleSetWithAttributeSelectorSuffixMatch(t *testing.T) {
	var stmts = RunParserTest(`[href$=pdf] { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserEmptyRuleSetWithTypeSelectorGroup(t *testing.T) {
	var stmts = RunParserTest(`div, span, html { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserEmptyRuleSetWithComplexSelector(t *testing.T) {
	var stmts = RunParserTest(`div#myId.first-name.last-name, span, html, .first-name, .last-name { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserNestedRuleSetSimple(t *testing.T) {
	var stmts = RunParserTest(`div, span, html { .foo { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserNestedRuleSetSimple2(t *testing.T) {
	var stmts = RunParserTest(`div, span, html { .foo { color: red; background: blue; } text-align: text; float: left; }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserNestedRuleWithParentSelector(t *testing.T) {
	var stmts = RunParserTest(`div, span, html { & { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserPropertyNameBorderWidth(t *testing.T) {
	var stmts = RunParserTest(`div { border-width: 3px 3px 3px 3px; }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserNestedProperty(t *testing.T) {
	var stmts = RunParserTest(`div {
		border: {
			width: 3px;
			color: #000;
		}
	}`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserPropertyNameBorderWidthInterpolation(t *testing.T) {
	var stmts = RunParserTest(`div { border-#{ $width }: 3px 3px 3px 3px; }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserPropertyNameBorderWidthInterpolation2(t *testing.T) {
	var stmts = RunParserTest(`div { #{ $name }: 3px 3px 3px 3px; }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserPropertyNameBorderWidthInterpolation3(t *testing.T) {
	var stmts = RunParserTest(`div { #{ $name }-left: 3px; }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserImportRuleWithUnquoteUrl(t *testing.T) {
	var stmts = RunParserTest(`@import url(../foo.css);`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserImportRuleWithUrl(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`@import url("http://foo.com/bar.css");`)
	assert.Equal(t, 1, len(*stmts))

	rule, ok := (*stmts)[0].(*ast.ImportStmt)
	assert.True(t, ok, "Convert to ImportStmt OK")
	assert.NotNil(t, rule)
}

func TestParserImportRuleWithString(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`@import "foo.css";`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserImportRuleWithMedia(t *testing.T) {
	var stmts = RunParserTest(`@import url("foo.css") screen;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserImportRuleWithMultipleMediaTypes(t *testing.T) {
	var stmts = RunParserTest(`@import url("bluish.css") projection, tv;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserImportRuleWithMediaTypeAndColorFeature(t *testing.T) {
	var stmts = RunParserTest(`@import url(color.css) screen and (color);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserImportRuleWithMediaTypeAndMaxWidthFeature(t *testing.T) {
	var stmts = RunParserTest(`@import url(color.css) screen and (max-width: 300px);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserImportRuleWithMedia2(t *testing.T) {
	var stmts = RunParserTest(`@import url("foo.css") screen and (orientation:landscape);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQuerySimple(t *testing.T) {
	var stmts = RunParserTest(`@media screen { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryNotScreen(t *testing.T) {
	var stmts = RunParserTest(`@media not screen { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryOnlyScreen(t *testing.T) {
	var stmts = RunParserTest(`@media only screen { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryAllAndMinWidth(t *testing.T) {
	var stmts = RunParserTest(`@media all and (min-width:500px) {  .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryMinWidth(t *testing.T) {
	var stmts = RunParserTest(`@media (min-width:500px) {  .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryOrientationPortrait(t *testing.T) {
	var stmts = RunParserTest(`@media (orientation: portrait) { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryMultipleWithComma(t *testing.T) {
	var stmts = RunParserTest(`@media screen and (color), projection and (color) { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryNone(t *testing.T) {
	var stmts = RunParserTest(`@media { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryNotAndMonoChrome(t *testing.T) {
	var stmts = RunParserTest(`@media not all and (monochrome) { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryJustAll(t *testing.T) {
	var stmts = RunParserTest(`@media all { .red { color: red; } }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryWithExpr1(t *testing.T) {
	var code = `
@media #{$media} {
  .sidebar {
    width: 500px;
  }
}
	`
	var stmts = RunParserTest(code)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryWithExpr2(t *testing.T) {
	var code = `
@media #{$media} and ($feature: $value) {
  .sidebar {
    width: 500px;
  }
}
	`
	var stmts = RunParserTest(code)
	assert.Equal(t, 1, len(*stmts))
}

/*
func TestParserMediaQueryNestedInRuleSet(t *testing.T) {
	var code = `
h6, .h6 {
  margin: 0 0 10px;
  line-height: 20px;
  font-size: 12px; }
  @media screen and (min-width: 960px) {
    h6, .h6 {
      font-size: 13px;
      margin: 0 0 15px;
  }
}
	`
	var stmts = RunParserTest(code)
	assert.Equal(t, 1, len(*stmts))
}
*/

func TestParserMediaQueryWithVendorPrefixFeature(t *testing.T) {
	// FIXME: 'min--moz-device-pixel-ratio' will become '-moz-device-pixel-ratio'
	var stmts = RunParserTest(`@media (-webkit-min-device-pixel-ratio: 2), (min--moz-device-pixel-ratio: 2) {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserMediaQueryNested(t *testing.T) {
	var code = `
@media screen {
  .sidebar {
    @media (orientation: landscape) {
      width: 500px;
    }
  }
}
	`
	var stmts = RunParserTest(code)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfTrueStmt(t *testing.T) {
	var stmts = RunParserTest(`@if true {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfFalseElseStmt(t *testing.T) {
	var stmts = RunParserTest(`@if false {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfFalseOrTrueElseStmt(t *testing.T) {
	var stmts = RunParserTest(`@if false or true {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfTrueAndTrueOrFalseElseStmt(t *testing.T) {
	var stmts = RunParserTest(`@if true and true or true {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfTrueAndTrueOrFalseElseStmt2(t *testing.T) {
	var stmts = RunParserTest(`@if (true and true) or true {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonGreaterThan(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) > 2 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonGreaterEqual(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) >= 2 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonLessThan(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) < 2 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonLessEqual(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) <= 2 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonEqual(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) == 6 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonUnequal(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) != 6 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonUnequalElseIf(t *testing.T) {
	var stmts = RunParserTest(`@if (3+3) != 6 {  } @else if (3+3) == 6 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfComparisonAndLogicalExpr(t *testing.T) {
	var stmts = RunParserTest(`@if 3 > 1 and 4 < 10 and 5 > 3 {  } @else if (3+3) == 6 {  } @else {  }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserIfDeclBlock(t *testing.T) {
	RunParserTest(`
@if $i == 1 {
	color: #111;
} @else if $i == 2 {
	color: #222;
} @else if $i == 3 {
	color: #333;
} @else {
	color: red;
	background: url(../background.png);
}
	`)
}

func TestParserForStmtSimple(t *testing.T) {
	var stmts = RunParserTest(`@for $var from 1 through 20 { }`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserForStmtExprReduce(t *testing.T) {
	var stmts = RunParserTest(`@for $var from 2 * 3 through 20 * 5 + 10 { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserForStmtRangeOperator(t *testing.T) {
	var stmts = RunParserTest(`@for $var in 1 .. 10 { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserForStmtRangeOperatorWithExpr(t *testing.T) {
	var stmts = RunParserTest(`@for $var in 2 + 3 .. 10 * 10 { }`)
	assert.Equal(t, 1, len(*stmts))

}

func TestParserWhileStmt(t *testing.T) {
	code := `
$i: 6;
@while $i > 0 { $i: $i - 2; }
`
	var stmts = RunParserTest(code)
	assert.Equal(t, 2, len(*stmts))
}

func TestParserCSS3Gradient(t *testing.T) {
	// some test cases from htmldog
	// @see http://www.htmldog.com/guides/css/advanced/gradients/
	var buffers = []string{
		`div { background: repeating-linear-gradient(white, black 10px, white 20px); }`,
		`div { background: linear-gradient(135deg, hsl(36,100%,50%) 10%, hsl(72,100%,50%) 60%, white 90%); }`,
		`div { background: linear-gradient(black 0, white 100%); }`,
		`div { background: radial-gradient(#06c 0, #fc0 50%, #039 100%); }`,
		`div { background: linear-gradient(red 0%, green 33.3%, blue 66.7%, black 100%); }`,
		`div { background: -webkit-radial-gradient(100px 200px, circle closest-side, black, white); }`,
	}
	for _, buffer := range buffers {
		var stmts = RunParserTest(buffer)
		assert.Equal(t, 1, len(*stmts))
	}
}

func TestParserPropertyListExpr(t *testing.T) {
	var buffers []string = []string{
		`div { width: 1px; }`,
		`div { width: 2px 3px; }`,
		`div { width: 4px, 5px, 6px, 7px; }`,
		`div { width: 4px, 5px 6px, 7px; }`,
		`div { width: 10px 3px + 7px 20px; }`,
		// `div { width: 10px, 3px + 7px, 20px; }`,
	}
	for _, buffer := range buffers {
		var stmts = RunParserTest(buffer)
		assert.Equal(t, 1, len(*stmts))
	}
}

func TestParserFontCssSlash(t *testing.T) {
	// should be plain CSS, no division
	// TODO: verify this case
	var stmts = RunParserTest(`.foo { font: 12px/24px; }`)
	t.Logf("%+v\n", stmts)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtWithBooleanTrue(t *testing.T) {
	var block = RunParserTest(`$foo: true;`)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithBooleanFalse(t *testing.T) {
	var block = RunParserTest(`$foo: false;`)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithNull(t *testing.T) {
	var stmts = RunParserTest(`$foo: null;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtList(t *testing.T) {
	var stmts = RunParserTest(`$foo: 1 2 3 4;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtListWithParenthesis(t *testing.T) {
	var stmts = RunParserTest(`$foo: (1 2 3 4);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtMap(t *testing.T) {
	var stmts = RunParserTest(`$foo: (bar: 1, foo: 2);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtCommaSepList(t *testing.T) {
	var stmts = RunParserTest(`$foo: (1,2,3,4);`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtWithMorePlus(t *testing.T) {
	var stmts = RunParserTest(`$foo: 12px + 20px + 20px;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtWithExprDefaultFlag(t *testing.T) {
	var stmts = RunParserTest(`$foo: 12px + 20px + 20px !default;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtWithExprOptionalFlag(t *testing.T) {
	var block = RunParserTest(`$foo: 12px + 20px + 20px !optional;`)
	t.Logf("%+v\n", block)
}

func TestParserAssignStmtWithComplexExpr(t *testing.T) {
	var stmts = RunParserTest(`$foo: 12px * (20px + 20px) + 4px / 2;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtWithInterpolation(t *testing.T) {
	var stmts = RunParserTest(`$foo: #{ 10 + 20 }px;`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserAssignStmtLengthPlusLength(t *testing.T) {
	var stmts = RunParserTest(`$foo: 10px + 20px;`)
	assert.Equal(t, 1, len(*stmts))
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtNumberPlusNumberMulLength(t *testing.T) {
	var stmts = RunParserTest(`$foo: (10 + 20) * 3px;`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorAddOperation(t *testing.T) {
	var stmts = RunParserTest(`$foo: #000 + 10;`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorMulOperation(t *testing.T) {
	var stmts = RunParserTest(`$foo: #010101 * 20;`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithHexColorDivOperation(t *testing.T) {
	var stmts = RunParserTest(`$foo: #121212 / 2;`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithPxValue(t *testing.T) {
	var stmts = RunParserTest(`$foo: 10px;`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithSolveableVariableRef(t *testing.T) {
	var stmts = RunParserTest(`
	$a: 10px; 
	$b: 10px;
	$c: $a + $b;
	`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithUnknownVariableRef(t *testing.T) {
	var stmts = RunParserTest(`
	$a: 10px; 
	$b: 10px;
	$c: 3 * ($a + $b) + $c;
	`)
	t.Logf("%+v\n", stmts) // should be 60px + $c
}

func TestParserAssignStmtWithFunctionCall(t *testing.T) {
	var stmts = RunParserTest(`$foo: go();`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallIntegerArgument(t *testing.T) {
	var stmts = RunParserTest(`$foo: go(1,2,3);`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallFunctionCallArgument(t *testing.T) {
	var stmts = RunParserTest(`$foo: go(bar());`)
	t.Logf("%+v\n", stmts)
}

func TestParserAssignStmtWithFunctionCallVariableArgument(t *testing.T) {
	var stmts = RunParserTest(`$foo: go($a,$b,$c);`)
	t.Logf("%+v\n", stmts)
}

func TestParserMixinSimple(t *testing.T) {
	RunParserTest(`
@mixin silly-links {
  a {
    color: blue;
    background-color: red;
  }
}
	`)
}

func TestParserMixinArguments(t *testing.T) {
	RunParserTest(`
@mixin colors($text, $background, $border) {
  color: $text;
  background-color: $background;
  border-color: $border;
}
	`)
}

func TestParserMixinContentDirective(t *testing.T) {
	RunParserTest(`
@mixin apply-to-ie6-only {
  * html {
    @content;
  }
}
	`)
}

func TestParserExtendClassSelector(t *testing.T) {
	RunParserTest(`@extend .foo-bar;`)
}

func TestParserExtendIdSelector(t *testing.T) {
	RunParserTest(`@extend #myId;`)
}

func TestParserExtendComplexSelector(t *testing.T) {
	RunParserTest(`@extend #myId > .foo-bar;`)
}

func TestParserInclude(t *testing.T) {
	RunParserTest(`
		@include apply-to-ie6-only;
	`)
}

func TestParserIncludeWithContentBlock(t *testing.T) {
	var stmts = RunParserTest(`
		@include apply-to-ie6-only {
			color: white;
		};
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionSimple(t *testing.T) {
	var stmts = RunParserTest(`
@function grid-width($n) {
  @return $n * $grid-width + ($n - 1) * $gutter-width;
}
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionSimple2(t *testing.T) {
	var stmts = RunParserTest(`
@function exists($name) {
  @return variable-exists($name);
}
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionSimple3(t *testing.T) {
	var stmts = RunParserTest(`
@function f() { }
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionSimple4(t *testing.T) {
	var stmts = RunParserTest(`
@function f() {
  $foo: hi;
  @return g();
}
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionSimple5(t *testing.T) {
	var stmts = RunParserTest(`
@function g() {
  @return variable-exists(foo);
}
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionWithAssignments(t *testing.T) {
	var stmts = RunParserTest(`
@function g() {
  $a: 2 * 10;
  @return $a * 99;
}
	`)
	assert.Equal(t, 1, len(*stmts))
}

func TestParserFunctionCallKeywordArguments(t *testing.T) {
	var stmts = RunParserTest(`
	@function foo($a, $b) {
		@return $a + $b;
	}
	$c: foo($b: 2, $a: 1);
	`)
	assert.Equal(t, 2, len(*stmts))
}

func TestParserMassiveRules(t *testing.T) {
	var buffers []string = []string{
		`div { width: auto; }`,
		`div { width: 100px }`,
		`div { width: 100pt }`,
		`div { width: 100em }`,
		`div { width: 100rem }`,
		`div { padding: 10px 10px; }`,
		`div { padding: 10px 10px 20px 30px; }`,
		`div { padding: 10px + 10px; }`,
		`div { padding: 10px + 10px * 3; }`,
		`div { color: red; }`,
		`div { color: rgb(255,255,255); }`,
		`div { color: rgba(255,255,255,0); }`,
		`div { background-image: url("../images/foo.png"); }`,
		// `div { color: #ccddee; }`,
	}
	for _, buffer := range buffers {
		t.Logf("%s\n", buffer)
		var p = NewParser(runtime.NewContext())
		var stmts = p.ParseScss(buffer)
		t.Logf("%+v\n", stmts)
	}
}

/*
func TestParserIfStmtTrueCondition(t *testing.T) {
	p := NewParser(runtime.NewContext())
	block := p.ParseScss(`
	div {
		@if true {
			color: red;
		}
	}
	`)
	_ = block
}
*/

func TestParserTypeSelector(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`div { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserClassSelector(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`.foo-bar { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserDescendantCombinatorSelector(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`
	.foo
	.bar
	.zoo { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserAdjacentCombinator(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`.foo + .bar { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserGeneralSiblingCombinator(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`.foo ~ .bar { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func TestParserChildCombinator(t *testing.T) {
	p := NewParser(runtime.NewContext())
	stmts := p.ParseScss(`.foo > .bar { width: auto; }`)
	ruleset, ok := (*stmts)[0].(*ast.RuleSet)
	assert.True(t, ok)
	assert.NotNil(t, ruleset)
}

func BenchmarkParserClassSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(runtime.NewContext())
		p.ParseScss(`.foo-bar {}`)
	}
}

func BenchmarkParserAttributeSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(runtime.NewContext())
		p.ParseScss(`input[type=text] {  }`)
	}
}

func BenchmarkParserComplexSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(runtime.NewContext())
		p.ParseScss(`div#myId.first-name.last-name, span, html, .first-name, .last-name { }`)
	}
}

func BenchmarkParserMediaQueryAllAndMinWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(runtime.NewContext())
		p.ParseScss(`@media all and (min-width:500px) {  .red { color: red; } }`)
	}
}

func BenchmarkParserOverAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p = NewParser(runtime.NewContext())
		p.ParseScss(`div#myId.first-name.last-name {
			.foo-bar {
				color: red;
				background: #fff;
				border-radius: 10px;
			}

			@for $i from 1 through 100 { }
			@if $i == 1 {
				color: #111;
			} @else if $i == 2 {
				color: #222;
			} @else if $i == 3 {
				color: #333;
			} @else {
				color: red;
				background: url(../background.png);
			}

			div { width: auto; }
			div { width: 100px }
			div { width: 100pt }
			div { width: 100em }
			div { width: 100rem }
			div { padding: 10px 10px; }
			div { padding: 10px 10px 20px 30px; }
			div { padding: 10px + 10px; }
			div { padding: 10px + 10px * 3; }
			div { color: red; }
			div { color: rgb(255,255,255); }
			div { color: rgba(255,255,255,0); }
			div { background-image: url("../images/foo.png"); }
		}`)
	}

}
