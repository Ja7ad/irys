install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest


check:
	gofumpt -l -w .
	golangci-lint run --timeout=20m0s