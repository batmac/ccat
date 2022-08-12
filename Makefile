GIT=$(shell git tag|tail -n1)
VERSION="post-${GIT}-dev"

all: ccat readme tests

ccat: */*/*.go */*/*/*.go go.mod go.sum
	./build.sh

readme: ccat
	cp README.header.md README.md
	echo >> README.md
	echo '```'      >> README.md
	sh -c './ccat --fullhelp 2>&1'  >> README.md
	echo '```'      >> README.md

tests: ccat
	go test -v ./...
	scripts/test_compression_e2e.sh testdata/compression/


janitor:
	golangci-lint --go=1.19 run --disable-all -E misspell --fix ./...
	golangci-lint --go=1.19 run ./...
	gofumpt -w -l .
	gosec -severity high ./...
	govulncheck ./...
	go list -json -deps ./... | nancy sleuth
	pre-commit autoupdate

docker-local: tests
	-rm ccat
	docker build --compress -t batmac/ccat:${VERSION} .

thanks:
	gothanks
