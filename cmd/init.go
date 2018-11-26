package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"

	"github.com/marhaupe/json-to-struct/internal"
)

func Start() {
	data, err := ioutil.ReadFile("./cmd/testdata.json")
	if err != nil {
		panic(err)
	}
	jsonString := string(data)
	c := make(chan json.Token)
	n := make(chan ds.JSONNode)
	var wg sync.WaitGroup
	wg.Add(2)

	go internal.Lex(jsonString, c, &wg)
	go internal.Parse(n, c, &wg)

	wg.Wait()
	fmt.Println(<-n)
}
