gen:
	go get -u github.com/clipperhouse/gen
	go get -u github.com/clipperhouse/typewriter
	go get -u github.com/clipperhouse/stringer
	go generate github.com/c9s/c6/...
	go build github.com/c9s/c6/...
	gen add github.com/c9s/c6/typewriter/symtable

get:
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/stringer
	go get golang.org/x/text/encoding
	go get golang.org/x/text/transform
	go get github.com/stretchr/testify/assert
	go get golang.org/x/tools/cmd/goimports
	go get github.com/ajstarks/svgo/benchviz
	go get github.com/spf13/cobra
	go get github.com/axw/gocov/gocov

vet:
	go vet github.com/c9s/c6/...

lint:
	go get -u github.com/golang/lint/golint
	golint github.com/c9s/c6/...


test: get gen 
	go test github.com/c9s/c6/...


benchupdatebase:
	go test -run=NONE -bench=. github.com/c9s/c6/... >| benchmarks/old.txt

benchrecord:
	go test -run=NONE -bench=. github.com/c9s/c6/... >| benchmarks/new.txt

bench:
	go test -run=NONE -bench=. -benchmem github.com/c9s/c6/...

benchcmp: benchrecord
	benchcmp benchmarks/old.txt benchmarks/new.txt

benchviz: benchrecord
	benchcmp benchmarks/old.txt benchmarks/new.txt | benchviz -top=5 -left=5 > benchmarks/summary.svg

cross-toolchain:
	gox -build-toolchain

cross-compile:
	gox -output "build/{{.Dir}}.{{.OS}}_{{.Arch}}" github.com/c9s/c6/...

cover:
	go test -cover -coverprofile c6.cov -coverpkg github.com/c9s/c6/ast,github.com/c9s/c6/runtime,github.com/c9s/c6/parser,github.com/c9s/c6/compiler github.com/c9s/c6/compiler
