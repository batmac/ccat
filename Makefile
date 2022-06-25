GIT=$(shell git tag|tail -n1)
VERSION="post-${GIT}-dev"

all: ccat readme tests

ccat: *.go */*.go */*.go go.mod go.sum
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
	golangci-lint run --disable-all -E misspell --fix ./...
	gofumpt -w -l .
	gosec -severity high ./...
	golangci-lint run ./...
	govulncheck ./...
	go list -json -deps ./... | nancy sleuth --skip-update-check
	go test -cover -coverprofile coverage.out ./...
	echo gocovsh --profile coverage.out

docker-local: tests
	-rm ccat
	docker build --compress -t batmac/ccat:${VERSION} .

docker: tests
	-rm ccat
	docker buildx build --compress -t batmac/ccat:${VERSION} --platform=linux/arm,linux/amd64,linux/arm64 -f Dockerfile . --push

docker-release: tests
	-rm ccat
	docker buildx build --compress -t batmac/ccat:latest -t batmac/ccat:${VERSION} --platform=linux/arm,linux/amd64,linux/arm64 -f Dockerfile . --push

thanks:
	gothanks
