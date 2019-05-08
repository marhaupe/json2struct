output ?= "json2struct"
version := $(shell git describe --tags)

.PHONY: build install test testrace coverage

build:
	go build -o ${output} -ldflags="-X github.com/marhaupe/json2struct/cmd.version=${version}"

install: 
	go install 

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