all:
	go generate c6 c6/ast
	go install -x c6 c6/ast

test: all
	go test c6/ast
	go test c6

benchupdate:
	go test -run=NONE -bench=. c6 >| benchmarks/old.txt

benchcmp: all
	go test -run=NONE -bench=. c6 >| benchmarks/new.txt
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt

cov:
	go test -coverprofile=c6.cov c6
