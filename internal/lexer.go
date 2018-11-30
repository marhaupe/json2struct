package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Lex consumes a string s containing the JSON object and writes each
// JSON Token to c
func Lex(s string, c chan json.Token, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			if err := fmt.Errorf("Error: %v, exiting", r); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()
	defer wg.Done()
	defer close(c)
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
			log.Fatal(err)
		}
		c <- t
	}
}
