GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.7.0
GOSEC := github.com/securego/gosec/v2/cmd/gosec@v2.28.0
GOFUMPT := mvdan.cc/gofumpt@v0.11.0

build:
	go run magefiles/mage.go

test:
	go run magefiles/mage.go
	go run magefiles/mage.go buildfull test
	go run magefiles/mage.go buildminimal test

thanks:
	gothanks

# golangci-lint is expected in $PATH: it is still pinned to the v1 config
# schema, see https://github.com/batmac/ccat/issues/1164
janitor:
	golangci-lint run --disable-all -E misspell --fix ./...
	golangci-lint run ./...
	go run $(GOFUMPT) -w -l .
	go run $(GOSEC) -severity high ./...
	go run $(GOVULNCHECK) ./...

release:
	goreleaser release --clean
	echo "go to https://github.com/batmac/ccat/releases and create a new release"

macsign:
	gon .gon.hcl

.PHONY: build test thanks janitor release all clean macsign
