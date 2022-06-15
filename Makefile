all: ccat readme

ccat: *.go */*.go go.mod go.sum
	./build.sh

readme: ccat
	cp README.header.md README.md
	echo '```'      >> README.md
	sh -c './ccat -h 2>&1'  >> README.md
	echo '```'      >> README.md

tests: ccat
	go test -v ./...
	scripts/test_compression_e2e.sh testdata/compression/
