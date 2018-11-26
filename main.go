package main

import (
	"fmt"
	"io/ioutil"

	"github.com/marhaupe/json-to-struct/cmd"
)

func main() {
	data, err := ioutil.ReadFile("./cmd/testdata.json")
	if err != nil {
		panic(err)
	}
	jsonString := string(data)
	s := cmd.Generate(jsonString)
	fmt.Print(s)
}
