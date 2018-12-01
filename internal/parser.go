package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/marhaupe/json2struct/internal/ds"
)

type Parser struct {
	rootEl ds.JSONNode
	c      chan json.Token
	wg     *sync.WaitGroup
}

// Parse parsed JSON Tokens received from chan c
func Parse(n chan ds.JSONNode, c chan json.Token, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			if err := fmt.Errorf("Error: %v, exiting", r); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	}()
	p := Parser{c: c, wg: wg}
	p.parse()
	n <- p.rootEl
}

func (p *Parser) parse() {
	defer p.wg.Done()

	// Consuming first Token, { or [
	t := <-p.c

	// Parsing root delim accordingly
	switch t {
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
	for t := range p.c {
		if key == "" {
			key = fmt.Sprint(t)

			// By making every second token the key, } and ] will falsely be recognized
			// as keys instead of values
			if key == "}" || key == "]" {
				return obj
			}
		} else {
			done := p.buildUpElement(obj, key, t, json.Delim('}'))
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

	for t := range p.c {
		key := generateArrayKeyForToken(t)
		done := p.buildUpElement(arr, key, t, json.Delim(']'))
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
		}
	}
	panic(fmt.Sprintf("Error generating key for token %v", t))
}

func (p *Parser) parsePrimitive(key string, value json.Token) *ds.JSONPrimitive {
	prim := &ds.JSONPrimitive{Key: key}
	prim.Ptype = detectPrimitiveType(value)
	return prim
}

func detectPrimitiveType(t json.Token) ds.PrimitiveType {
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
	default:
		fmt.Printf("Could not determine datatype of field with value %v\n", t)
		return ds.String
	}
}
