package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/marhaupe/json2struct/internal/ds"
)

// Parser contains fields to enable parsing of JSON tokens to
// JSONElements
type Parser struct {
	rootEl ds.JSONNode
	lr     chan LexResult
	wg     *sync.WaitGroup
}

// ParseResult is meant to be used as a chan struct to enable sharing
// error messages between goroutines
type ParseResult struct {
	Node  ds.JSONNode
	Error error
}

// Parse parses JSON tokens received from chan lr and writes either the resulting
// error or the parsed JSONNode to chan pr
func Parse(pr chan ParseResult, lr chan LexResult, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			pr <- ParseResult{Node: nil, Error: err}
			close(pr)
		}
	}()

	p := Parser{lr: lr, wg: wg}
	p.parse()
	pr <- ParseResult{Node: p.rootEl, Error: nil}
	close(pr)
}

func (p *Parser) parse() {
	defer p.wg.Done()

	// Consuming first Token, { or [
	r := <-p.lr
	if r.Error != nil {
		panic(r.Error)
	}
	// Parsing root delim accordingly
	switch r.Token {
	case json.Delim('{'):
		p.buildRootObject()
	case json.Delim('['):
		p.buildRootArray()
	}
}

func (p *Parser) buildRootObject() {
	obj := p.parseObject("ROOT_OBJECT")
	obj.Root = true
	p.rootEl = obj
}

func (p *Parser) buildRootArray() {
	arr := p.parseArray("ROOT_ARRAY")
	arr.Root = true
	p.rootEl = arr
}

func (p *Parser) parseObject(objKey string) *ds.JSONObject {
	obj := &ds.JSONObject{Key: objKey}
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

func (p *Parser) parseArray(arrKey string) *ds.JSONArray {
	arr := &ds.JSONArray{Key: arrKey}

	for t := range p.lr {
		if t.Error != nil {
			panic(t.Error)
		}
		key := generateArrayKeyForToken(t.Token)
		done := p.buildUpElement(arr, key, t.Token, json.Delim(']'))
		if done {
			return arr
		}
	}
	return arr
}

func (p *Parser) buildUpElement(node ds.JSONNode, key string, t json.Token, closing json.Delim) bool {
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

func generateArrayKeyForToken(t json.Token) string {
	switch t {
	// These delims will get detected as keys
	case json.Delim('}'), json.Delim(']'):
		return "closing_delim"
	case json.Delim('{'):
		return "object_in_array"
	case json.Delim('['):
		return "array_in_array"
	default:
		switch detectPrimitiveType(t) {
		case ds.String:
			return "string_in_array"
		case ds.Bool:
			return "bool_in_array"
		case ds.Float:
			return "float_in_array"
		case ds.Int:
			return "int_in_array"
		case ds.Null:
			return "null_in_array"
		}
	}
	panic(fmt.Sprintf("Error generating key for token %v", t))
}

func (p *Parser) parsePrimitive(key string, value json.Token) *ds.JSONPrimitive {
	prim := &ds.JSONPrimitive{Key: key}
	prim.Datatype = detectPrimitiveType(value)
	return prim
}

func detectPrimitiveType(t json.Token) ds.Datatype {
	switch t.(type) {
	case float64:
		if t == math.Trunc(t.(float64)) {
			return ds.Int
		}
		return ds.Float
	case string:
		return ds.String
	case bool:
		return ds.Bool
	case nil:
		return ds.Null
	default:
		fmt.Printf("Could not determine datatype of field with value %v\n", t)
		return ds.String
	}
}
