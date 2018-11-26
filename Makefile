OUTPUT ?= "./json-to-struct"

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(OUTPUT)

install: 
	CGO_ENABLED=0 go install -a -installsuffix cgo

test:
	go test ./...

testCoverage:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html 