// Package generate contains fields to generate a Go struct from a given JSON string
package generate

import (
	"fmt"
	"sync"

	"github.com/marhaupe/json2struct/internal/lex"
	"github.com/marhaupe/json2struct/internal/parse"
)

// Generate takes a JSON string s and parses it accordingly
// The resulting string is generated code for a Go struct that e.g. can be used to
// marshal and unmarshal JSON
func Generate(s string) (string, error) {
	lr := make(chan lex.Result)
	pr := make(chan parse.Result)
	var wg sync.WaitGroup
	wg.Add(2)

	go lex.Lex(s, lr, &wg)
	go parse.Parse(pr, lr, &wg)
	wg.Wait()
	res := <-pr
	if res.Error != nil {
		return "", res.Error
	}
	return fmt.Sprint(res.Node), nil
}
