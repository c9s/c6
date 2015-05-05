all:
	go generate c6 c6/ast
	go install -x c6 c6/ast

test: all
	go test c6/ast
	go test c6

bench: all
	go test -bench=. c6


cov:
	go test -coverprofile=c6.cov c6
