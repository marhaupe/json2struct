package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

// LexResult is meant to be used as a chan struct to enable sharing
// error messages between goroutines
type LexResult struct {
	Token json.Token
	Error error
}

// Lex consumes a string s containing the JSON object and writes each
// result containing either the token or the error to c
func Lex(s string, lr chan LexResult, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			lr <- LexResult{Error: err, Token: nil}
			close(lr)
		}
	}()
	defer wg.Done()
	ok := json.Valid([]byte(s))
	if !ok {
		panic("Invalid json")
	}

	d := json.NewDecoder(strings.NewReader(s))
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			lr <- LexResult{Token: nil, Error: err}
		}
		lr <- LexResult{Token: t, Error: nil}
	}
	close(lr)
}
