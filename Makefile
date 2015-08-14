gen:
	go get -u github.com/clipperhouse/gen
	go get -u github.com/clipperhouse/typewriter
	go get -u github.com/clipperhouse/stringer
	go generate github.com/c9s/c6/...
	go build github.com/c9s/c6/...
	gen add github.com/c9s/c6/typewriter/symtable


vet:
	go vet github.com/c9s/c6/...


test: gen vet 
	go test github.com/c9s/c6/...


benchupdatebase:
	go test -run=NONE -bench=. github.com/c9s/c6/... >| benchmarks/old.txt

benchrecord:
	go test -run=NONE -bench=. github.com/c9s/c6/... >| benchmarks/new.txt

bench:
	go test -run=NONE -bench=. -benchmem github.com/c9s/c6/...

benchcmp: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt

benchviz: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt | benchviz -top=5 -left=5 > benchmarks/summary.svg

cross-toolchain:
	gox -build-toolchain

cross-compile:
	gox -output "build/{{.Dir}}.{{.OS}}_{{.Arch}}" c6/...

cover:
	go test -cover -coverprofile c6.cov -coverpkg github.com/c9s/c6/ast,github.com/c9s/c6/runtime,github.com/c9s/c6/parser c6/parser
