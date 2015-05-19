all:
	go generate c6/ast c6
	go build c6/ast c6/runtime c6

vendor:
	source goinstall

clean:
	go clean c6/...

install:
	go install c6/...

test: all
	go test c6/ast c6/runtime c6

benchupdatebase:
	go test -run=NONE -bench=. c6 >| benchmarks/old.txt

benchrecord:
	go test -run=NONE -bench=. c6 >| benchmarks/new.txt

bench:
	go test -run=NONE -bench=. -benchmem c6 

benchcmp: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt

benchviz: all benchrecord
	vendor/bin/benchcmp benchmarks/old.txt benchmarks/new.txt | benchviz -top=5 -left=5 > benchmarks/summary.svg

cross-toolchain:
	gox -build-toolchain

cross-compile:
	gox c6/...

cover:
	go test -cover -coverprofile c6.cov -coverpkg c6,c6/ast,c6/runtime c6

cover-annotate: cov
	vendor/bin/gocov convert c6.cov | vendor/bin/gocov annotate -
