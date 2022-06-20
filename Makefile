all: ccat readme tests

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


janitor:
	gosec -severity high ./...
	golangci-lint run ./...
	govulncheck ./...
	golangci-lint run --disable-all -E gofumpt -E misspell --fix ./...
	go test -cover -coverprofile coverage.out ./...
	echo gocovsh --profile coverage.out
