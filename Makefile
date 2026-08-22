GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.7.0
GOSEC := github.com/securego/gosec/v2/cmd/gosec@v2.28.0
GOFUMPT := mvdan.cc/gofumpt@v0.11.0
GOLANGCI_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.1

# advisories with no fixed version available, so govulncheck can never go green:
# GO-2026-5932 is x/crypto/openpgp, reached through go-selfupdate. See
# https://github.com/batmac/ccat/issues/1163
VULN_IGNORE := GO-2026-5932

build:
	go run magefiles/mage.go

test:
	go run magefiles/mage.go
	go run magefiles/mage.go buildfull test
	go run magefiles/mage.go buildminimal test

thanks:
	gothanks

# every tool runs via 'go run': the local toolchain builds it, so it always
# targets the same Go version as go.mod (a golangci-lint binary built with an
# older Go refuses to load a newer one)
janitor:
	go run $(GOLANGCI_LINT) run --enable-only misspell --fix ./...
	go run $(GOLANGCI_LINT) run ./...
	go run $(GOFUMPT) -w -l .
	go run $(GOSEC) -severity high -exclude-dir=old ./...
	@out="$$(go run $(GOVULNCHECK) ./... 2>&1 || true)"; \
	echo "$$out"; \
	if echo "$$out" | grep -E '^Vulnerability #' | grep -qv '$(VULN_IGNORE)'; then \
		echo "==> govulncheck found a vulnerability other than $(VULN_IGNORE)"; \
		exit 1; \
	fi

release:
	goreleaser release --clean
	echo "go to https://github.com/batmac/ccat/releases and create a new release"

macsign:
	gon .gon.hcl

.PHONY: build test thanks janitor release all clean macsign
