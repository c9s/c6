all:
	go generate github.com/c9s/c6/src/c6/ast github.com/c9s/c6/src/c6
	go build github.com/c9s/c6/src/c6/ast github.com/c9s/c6/src/c6

deps:
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/stringer
	go get github.com/stretchr/testify/assert
	go get golang.org/x/tools/cmd/goimports
	go get golang.org/x/tools/cmd/benchcmp
	go get github.com/ajstarks/svgo/benchviz
	go get github.com/axw/gocov/gocov

clean:
	go clean github.com/c9s/c6/src/c6/...

install:
	go install github.com/c9s/c6/src/c6/...

test: all
	go test github.com/c9s/c6/src/c6/ast github.com/c9s/c6/src/c6

benchupdatebase:
	go test -run=NONE -bench=. github.com/c9s/c6/src/c6 >| benchmarks/old.txt

benchrecord:
	go test -run=NONE -bench=. github.com/c9s/c6/src/c6 >| benchmarks/new.txt

bench:
	go test -run=NONE -bench=. -benchmem github.com/c9s/c6/src/c6

benchcmp: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt

benchviz: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt | benchviz -top=5 -left=5 > benchmarks/summary.svg

cross-toolchain:
	gox -build-toolchain

cross-compile:
	gox -output "build/{{.Dir}}.{{.OS}}_{{.Arch}}" github.com/c9s/c6/src/c6/...

cover:
	go test -cover -coverprofile c6.cov -coverpkg github.com/c9s/c6/src/c6,github.com/c9s/c6/src/c6/ast github.com/c9s/c6/src/c6

cover-annotate: cov
	vendor/bin/gocov convert c6.cov | vendor/bin/gocov annotate -
