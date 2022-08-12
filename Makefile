GIT=$(shell git tag|tail -n1)
VERSION="post-${GIT}-dev"

all:  readme tests

bootstrap:
	@go install github.com/magefile/mage@latest

janitor:
	golangci-lint --go=1.19 run --disable-all -E misspell --fix ./...
	golangci-lint --go=1.19 run ./...
	gofumpt -w -l .
	gosec -severity high ./...
	govulncheck ./...
	go list -json -deps ./... | nancy sleuth
	pre-commit autoupdate

thanks:
	gothanks
