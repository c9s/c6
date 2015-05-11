all:
	go generate c6/ast
	go generate c6
	go build c6/ast 
	go build c6/runtime
	go build c6

install:
	go install c6/...

test: all
	go test c6/ast
	go test c6/runtime
	go test c6

benchupdate:
	go test -run=NONE -bench=. c6 >| benchmarks/old.txt

benchcmp: all
	go test -run=NONE -bench=. c6 >| benchmarks/new.txt
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt

cov:
	go test -coverprofile=c6.cov c6
	go test -coverprofile=c6_ast.cov c6/ast
	go test -coverprofile=c6_runtime.cov c6/runtime

cov-annotate: cov
	vendor/bin/gocov convert c6.cov | vendor/bin/gocov annotate -
	vendor/bin/gocov convert c6_ast.cov | vendor/bin/gocov annotate -
	vendor/bin/gocov convert c6_runtime.cov | vendor/bin/gocov annotate -
