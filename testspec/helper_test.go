package testspec

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/c9s/c6/compiler"
	"github.com/c9s/c6/parser"
	"github.com/c9s/c6/runtime"
	"github.com/can3p/go-hrx/hrx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunSpecs(t *testing.T, hrxPath string, testFiles []string, testOnly string, ignoreInput func(string, string) bool) map[string][][]string {
	success := 0

	failedSpecs := map[string][][]string{}

	addFailure := func(spec, fname string, reason string) {
		if _, ok := failedSpecs[spec]; !ok {
			failedSpecs[spec] = [][]string{}
		}
		failedSpecs[spec] = append(failedSpecs[spec], []string{fname, reason})
	}

	reportSuccess := func(fname, input string) {
		if ignoreInput != nil && ignoreInput(fname, input) {
			assert.True(t, false, "[example %s] Input: %s - test passes, but is listed in the blacklisted object", fname, input)
			return
		}

		success++
	}

	testedCount := 0
	for _, fname := range testFiles {
		if testOnly != "" && testOnly != fname {
			continue
		}

		archive, err := hrx.OpenReader(path.Join(hrxPath, fname))
		require.NoErrorf(t, err, "[example %s] open hrx file", fname)

		inputFiles, err := getInputFiles(archive)

		require.NoErrorf(t, err, "[example %s] get input files", fname)

		for _, input := range inputFiles {
			if ignoreInput != nil && ignoreInput(fname, input) {
				continue
			}

			t.Logf("Processing Input file: %s - %s", fname, input)

			testedCount++

			if !assert.NotPanicsf(t, func() {
				var context = runtime.NewContext()
				var parser = parser.NewParser(context)

				baseName := path.Dir(input)
				errFname := path.Join(baseName, "error")
				warnFname := path.Join(baseName, "warning")
				outputFname := path.Join(baseName, "output.css")

				var expectedError string
				var expectedWarning string

				if _, err := fs.Stat(archive, errFname); err == nil {
					expected, err := fs.ReadFile(archive, errFname)

					require.NoErrorf(t, err, "[example %s] Input: %s", fname, errFname)
					expectedError = string(expected)
				}

				if _, warn := fs.Stat(archive, warnFname); warn == nil {
					expected, warn := fs.ReadFile(archive, warnFname)

					require.NoError(t, warn, "[example %s] Input: %s", fname, warnFname)
					expectedWarning = string(expected)
				}

				var stmts, parseErr = parser.ParseFile(archive, input)
				if parseErr != nil {
					if expectedError != "" {
						if !assert.Equalf(t, expectedError, parseErr.Error(), "[example %s] Input: %s", fname, errFname) {
							addFailure(fname, input, "parser_error_does_not_match")
						} else {
							reportSuccess(fname, input)
						}

						return
					}

					assert.NoErrorf(t, parseErr, "[example %s] Input: %s, parse failed", fname, input)
					addFailure(fname, input, "parse_failure")
					return
				}

				var b bytes.Buffer
				var warn bytes.Buffer

				wr := func(msg any) {
					// that's a silly implementation
					warn.WriteString(fmt.Sprintf("%s", msg))
				}

				var compiler = compiler.NewPrettyCompiler(&b, compiler.WithWarn(wr), compiler.WithDebug(wr))

				compileErr := compiler.Compile(stmts)

				if _, err := fs.Stat(archive, warnFname); err == nil {
					if !assert.Equal(t, expectedWarning, warn.String(), "[example %s] Input: %s - warning is expected", fname, input) {
						addFailure(fname, input, "unhandled_warning")
						return
					}
				}

				if expectedError != "" {
					if !assert.Errorf(t, compileErr, "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, "compiler_should_have_errored")
						return
					}

					if !assert.Equalf(t, expectedError, compileErr.Error(), "[example %s] Input: %s", fname, errFname) {
						addFailure(fname, input, "compiler_error_does_not_match")
						return
					}

					reportSuccess(fname, input)
				} else {
					if !assert.NoErrorf(t, compileErr, "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, fmt.Sprintf("compiler_unexpected_compile_error - %s", compileErr.Error()))
						return
					}

					expected, err := fs.ReadFile(archive, outputFname)
					if !assert.NoErrorf(t, err, "[example %s] Input: %s", fname, input) ||
						!assert.Equalf(t, string(expected), b.String(), "[example %s] Input: %s", fname, input) {
						addFailure(fname, input, "compiler_output_does_not_match")
						return
					}

					reportSuccess(fname, input)
				}
			}, "[example %s] Input: %s", fname, input) {
				addFailure(fname, input, "compiler_panic")
			}
		}
	}

	t.Logf("%d/%d spec files were successful", success, testedCount)
	return failedSpecs
}

func getHrxFiles(path string) ([]string, error) {
	names := []string{}

	err := fs.WalkDir(os.DirFS(path), ".", func(p string, d fs.DirEntry, err2 error) error {
		if strings.HasSuffix(d.Name(), ".hrx") {
			names = append(names, p)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return names, nil
}

func getInputFiles(fsys fs.FS) ([]string, error) {
	out := []string{}

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err2 error) error {
		if strings.HasSuffix(p, "input.scss") || strings.HasSuffix(p, "input.sass") {
			out = append(out, p)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
