OUTPUT?="json2struct"
VERSION=$(shell git describe --tags)
LDFLAGS="-X github.com/marhaupe/json2struct/cmd.version=$(VERSION)"

.PHONY: build install test testrace clean coverage 

build:
	go build -o $(OUTPUT) -ldflags=$(LDFLAGS)

install: 
	go install -ldflags=$(LDFLAGS)

test:
	go test ./...

testrace:
	go test -race ./...

clean:
	rm -rf dist

coverage: 
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html 
	rm cover.out