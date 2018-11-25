package internal

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"sync"
)

// Lex consumes a string s containing the JSON object and writes each
// JSON Token to c
func Lex(s string, c chan json.Token, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(c)
	ok := json.Valid([]byte(s))
	if !ok {
		panic("Invalid json!")
	}

	d := json.NewDecoder(strings.NewReader(s))
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		c <- t
	}
}
