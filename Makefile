get-dev:
	go get -u golang.org/x/tools/cmd/benchcmp
	go get -u golang.org/x/tools/cmd/stringer
	go get -u github.com/golang/lint/golint
	go get -u github.com/ajstarks/svgo/benchviz
	go get -t github.com/c9s/c6/...


test: 
	go test github.com/c9s/c6/...
	# The tests should pass aftr `go generete` too,
	# # This is to guard against changes that can break
	# # with `go generate`.
	go generate github.com/c9s/c6/...
	go test github.com/c9s/c6/...

vet: 
	go vet github.com/c9s/c6/...

lint:
	golint github.com/c9s/c6/...

cover:
	go test -cover -coverprofile c6.cov -coverpkg github.com/c9s/c6/ast,github.com/c9s/c6/runtime,github.com/c9s/c6/parser,github.com/c9s/c6/compiler github.com/c9s/c6/compiler

benchmark: 
	go test -run=NONE -bench=. github.com/c9s/c6/... >| benchmarks/new.txt
	benchcmp benchmarks/old.txt benchmarks/new.txt

benchviz: benchrecord
	benchcmp benchmarks/old.txt benchmarks/new.txt | benchviz -top=5 -left=5 > benchmarks/summary.svg

cross-toolchain:
	gox -build-toolchain

cross-compile:
	gox -output "build/{{.Dir}}.{{.OS}}_{{.Arch}}" github.com/c9s/c6/...
