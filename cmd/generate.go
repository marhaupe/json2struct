package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/marhaupe/json2struct/internal/ds"

	"github.com/marhaupe/json2struct/internal"
)

func Generate(s string) string {
	c := make(chan json.Token)
	n := make(chan ds.JSONNode)
	var wg sync.WaitGroup
	wg.Add(2)
	go internal.Lex(s, c, &wg)
	go internal.Parse(n, c, &wg)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("Error %v, exiting", r)
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()
	wg.Wait()

	return fmt.Sprint(<-n)
}
