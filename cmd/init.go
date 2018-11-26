package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"

	"github.com/marhaupe/json-to-struct/internal"
)

func Start() {
	j := `{ 
		"Hallo": "Hey", 
		"Schoen": true,
		"Number": 5999
		}`
	c := make(chan json.Token)
	n := make(chan ds.JSONNode)
	var wg sync.WaitGroup
	wg.Add(2)
	go internal.Lex(j, c, &wg)
	go internal.Parse(n, c, &wg)
	wg.Wait()
	fmt.Println(<-n)
}
