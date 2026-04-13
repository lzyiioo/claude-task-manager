.PHONY: build run install clean test

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY := ctm

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/main.go

run:
	go run ./cmd/main.go

install: build
	mv $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)

test:
	go test ./...

# Cross-compilation
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-linux-amd64 ./cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-darwin-amd64 ./cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-darwin-arm64 ./cmd/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-windows-amd64.exe ./cmd/main.go
