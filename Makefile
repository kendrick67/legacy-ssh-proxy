BINARY := legacy-ssh-proxy
PKG := ./...

.PHONY: build test lint test-wheel clean

build:
	go build -o bin/$(BINARY) ./cmd/legacy-ssh-proxy

test:
	go test -race $(PKG)

lint:
	go vet $(PKG)

test-wheel:
	./scripts/build_whl_for_local_testing.sh

clean:
	rm -rf bin dist python/src/legacy_ssh_proxy/bin python/dist python/dist-raw python/build ./tmp/test-venv
