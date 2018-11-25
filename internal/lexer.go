package internal

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"
)

func Lex(s string, c chan json.Token, wg *sync.WaitGroup) {
	l := &Lexer{}
	l.lex(s, c, wg)
}

type Lexer struct {
	root ds.JSONElement
}

func (l *Lexer) lex(s string, c chan json.Token, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(c)
	ok := json.Valid([]byte(s))
	if !ok {
		panic("Invalid json!")
	}
	d := json.NewDecoder(strings.NewReader(s))
	delim, err := d.Token()
	if err != nil {
		log.Fatal(err)
	}
	switch delim {
	case json.Delim('{'):
		l.root = &ds.JSONObject{
			Root: true,
		}
		c <- delim
	case json.Delim('['):
		l.root = &ds.JSONArray{
			Root: true,
		}
		c <- delim
	default:
		panic("Error parsing")
	}
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
