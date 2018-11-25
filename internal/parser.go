package internal

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"
)

type Parser struct {
	root ds.JSONElement
}

func Parse(c chan json.Token, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range c {
		fmt.Printf("%T: %v \n", t, t)
	}
}
