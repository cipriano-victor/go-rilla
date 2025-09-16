.PHONY: test race lint fmt


test:
	go test ./...


race:
	go test -race ./...


fmt:
	gofumpt -w . || gofmt -w .


lint:
	golangci-lint run || true