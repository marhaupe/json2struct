package parse

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/marhaupe/json2struct/internal/lex"
)

// parse contains fields to enable parsing of JSON tokens to
// JSONElements
type parse struct {
	rootEl JSONNode
	lr     chan lex.Result
	wg     *sync.WaitGroup
}

// Result is meant to be used as a chan struct to enable sharing
// error messages between goroutines
type Result struct {
	Node  JSONNode
	Error error
}

// Parse parses JSON tokens received from chan lr and writes either the resulting
// error or the parsed JSONNode to chan pr
func Parse(pr chan Result, lr chan lex.Result, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			pr <- Result{Node: nil, Error: err}
			close(pr)
		}
	}()

	p := parse{lr: lr, wg: wg}
	p.parse()
	pr <- Result{Node: p.rootEl, Error: nil}
	close(pr)
}

func (p *parse) parse() {
	defer p.wg.Done()

	// Consuming first Token, { or [
	r := <-p.lr
	if r.Error != nil {
		panic(r.Error)
	}
	// Parsing root delim accordingly
	switch r.Token {
	case json.Delim('{'):
		p.rootEl = p.parseObject("")
	case json.Delim('['):
		p.rootEl = p.parseArray("")
	}
}

func (p *parse) parseObject(objKey string) *JSONObject {
	obj := &JSONObject{Key: objKey}
	var key string
	for r := range p.lr {
		if r.Error != nil {
			panic(r.Error)
		}
		if key == "" {
			key = fmt.Sprint(r.Token)

			// By making every second token the key, } and ] will falsely be recognized
			// as keys instead of values
			if key == "}" || key == "]" {
				return obj
			}
		} else {
			done := p.buildUpElement(obj, key, r.Token, json.Delim('}'))
			if done {
				return obj
			}
			key = ""
		}
	}
	return obj
}

func (p *parse) parseArray(arrKey string) *JSONArray {
	arr := &JSONArray{Key: arrKey}

	for t := range p.lr {
		if t.Error != nil {
			panic(t.Error)
		}
		key := ""
		done := p.buildUpElement(arr, key, t.Token, json.Delim(']'))
		if done {
			return arr
		}
	}
	return arr
}

func (p *parse) buildUpElement(node JSONNode, key string, t json.Token, closing json.Delim) bool {
	switch t {
	case closing:
		return true
	case json.Delim('{'):
		node.AddChild(p.parseObject(key))
	case json.Delim('['):
		node.AddChild(p.parseArray(key))
	default:
		node.AddChild(p.parsePrimitive(key, t))
	}
	return false
}

func (p *parse) parsePrimitive(key string, value json.Token) *JSONPrimitive {
	prim := &JSONPrimitive{Key: key}
	prim.Datatype = detectPrimitiveType(value)
	return prim
}

func detectPrimitiveType(t json.Token) Datatype {
	switch t.(type) {
	case float64:
		if t == math.Trunc(t.(float64)) {
			return Int
		}
		return Float
	case string:
		return String
	case bool:
		return Bool
	case nil:
		return Null
	default:
		fmt.Printf("Could not determine datatype of field with value %v\n", t)
		return String
	}
}
