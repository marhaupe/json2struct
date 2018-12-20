package internal

import (
	"fmt"
	"sync"
)

func Generate(s string) (string, error) {
	c := make(chan LexResult)
	r := make(chan ParseResult)
	var wg sync.WaitGroup
	wg.Add(2)

	go Lex(s, c, &wg)
	go Parse(r, c, &wg)
	wg.Wait()
	res := <-r
	if res.Error != nil {
		return "", res.Error
	}
	return fmt.Sprint(res.Node), nil
}
