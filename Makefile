BINARY := legacy-ssh-proxy
PKG := ./...

.PHONY: build test lint clean

build:
	go build -o bin/$(BINARY) ./cmd/legacy-ssh-proxy

test:
	go test -race $(PKG)

lint:
	go vet $(PKG)

clean:
	rm -rf bin dist
