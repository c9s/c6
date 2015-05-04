all:
	go generate c6 c6/ast
	go build -x c6 c6/ast

test:
	go test -x c6
	go test -x c6/ast

cov:
	go test -coverprofile=c6.cov c6
