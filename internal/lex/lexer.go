// Package lex provides fields to enable the tokenization of a JSON, meant to be
// used in conjunction with package parse.
package lex

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

// Result is the result of a the lexing of JSON. It either contains Result, representing
// a json.Token like { or }, or and error
type Result struct {
	Token json.Token
	Error error
}

// Lex consumes a string s containing the JSON object and writes each
// result containing either the token or the error to lexRes
func Lex(s string, lexRes chan Result, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			lexRes <- Result{Error: err, Token: nil}
			close(lexRes)
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
			lexRes <- Result{Token: nil, Error: err}
		}
		lexRes <- Result{Token: t, Error: nil}
	}
	close(lexRes)
}
