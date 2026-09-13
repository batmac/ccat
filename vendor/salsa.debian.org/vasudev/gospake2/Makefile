PKG = "salsa.debian.org/vasudev/gospake2"
PKG_LIST := $(shell go list $(PKG)/...)

.PHONY: all build dep lint test clean cover

all: build

build: dep
	go build -i -v

cover: dep
	./scripts/coverage.sh

lint: dep
	go fmt ${PKG_LIST}
	golint -set_exit_status ${PKG_LIST}
	go vet ${PKG_LIST}

dep:
	./scripts/getdeps.sh

test: dep
	go test -race -short ${PKG_LIST}
