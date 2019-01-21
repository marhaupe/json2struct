package internal

import (
	"fmt"
	"sync"
)

// Generate takes a JSON string s and lexes and parses that string to JSONElements.
// The resulting string is generated code for a golang struct that can be used to
// marshal and unmarshal JSON
func Generate(s string) (string, error) {
	lr := make(chan LexResult)
	pr := make(chan ParseResult)
	var wg sync.WaitGroup
	wg.Add(2)

	go Lex(s, lr, &wg)
	go Parse(pr, lr, &wg)
	wg.Wait()
	res := <-pr
	if res.Error != nil {
		return "", res.Error
	}
	return fmt.Sprint(res.Node), nil
}
