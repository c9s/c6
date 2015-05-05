all:
	go generate c6 c6/ast
	go install -x c6 c6/ast

test: all
	go test -i c6/ast
	go test -i c6

bench: all
	go test -i -bench=. c6


cov:
	go test -coverprofile=c6.cov c6
