.PHONY: test race lint fmt


test:
	go test ./...

fmt:
	gofumpt -w . || gofmt -w .


lint:
	golangci-lint run || true