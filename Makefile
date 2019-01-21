OUTPUT ?= "./json2struct"

build:
	go build -o $(OUTPUT)

install: 
	go install 

test:
	go test ./...

testrace:
	go test -race ./...

coverage:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html 
	rm cover.out