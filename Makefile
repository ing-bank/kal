GOCDM=go
GOBUILD=$(GOCDM) build
GOMOD=$(GOCDM) mod
BINARY=kal
GOHOME=~/go

.PHONY: install
.PHONY: test
.PHONY: lint

all: clean tidy install

install:
	CGO_ENABLED=0 go install .

test: 
    CGO_ENABLED=0 go test -timeout=5m $(shell go list ./... | grep -v /vendor/)

lint: 
    golangci-lint run --enable misspell --out-format=colored-line-number --timeout 10m


