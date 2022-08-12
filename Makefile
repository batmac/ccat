
build:
	@mage

thanks:
	gothanks

janitor:
	golangci-lint --go=1.19 run --disable-all -E misspell --fix ./...
	golangci-lint --go=1.19 run ./...
	gofumpt -w -l .
	gosec -severity high ./...
	govulncheck ./...
	go list -json -deps ./... | nancy sleuth
	pre-commit autoupdate
