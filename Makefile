
build:
	go run magefiles/mage.go

test:
	go run magefiles/mage.go
	go run magefiles/mage.go buildfull test
	go run magefiles/mage.go buildminimal test

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
