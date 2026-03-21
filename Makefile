VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build dev test lint clean install

build:
	go build -ldflags "$(LDFLAGS)" -o bin/idops ./cmd/idops

dev:
	go run ./cmd/idops $(ARGS)

test:
	go test ./... -v

lint:
	golangci-lint run ./...

install: build
	cp bin/idops /usr/local/bin/idops

clean:
	rm -rf bin/
