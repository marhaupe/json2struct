package lex

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

// Result is meant to be used as a chan struct to enable sharing
// error messages between goroutines
type Result struct {
	Token json.Token
	Error error
}

// Lex consumes a string s containing the JSON object and writes each
// result containing either the token or the error to c
func Lex(s string, lr chan Result, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			lr <- Result{Error: err, Token: nil}
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
			lr <- Result{Token: nil, Error: err}
		}
		lr <- Result{Token: t, Error: nil}
	}
	close(lr)
}
