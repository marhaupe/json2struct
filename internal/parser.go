package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/marhaupe/json2struct/internal/ds"
)

const FILLER_KEY string = "FILLER_KEY"

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
	obj := p.parseObject(FILLER_KEY)
	obj.Root = true
	p.rootEl = obj
}

func (p *Parser) buildRootArray() {
	arr := p.parseArray(FILLER_KEY)
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
			switch t {
			case json.Delim('{'):
				obj.AddChild(p.parseObject(key))
			case json.Delim('['):
				obj.AddChild(p.parseArray(key))
			case json.Delim('}'):
				return obj
			default:
				obj.AddChild(p.parsePrimitive(key, t))
			}
			key = ""
		}
	}
	return obj
}

func (p *Parser) parseArray(arrKey string) *ds.JSONArray {
	arr := &ds.JSONArray{Key: arrKey}

	// Contents of array do not need a key; filling value with filler key
	var key string
	for t := range p.c {
		switch t {
		case json.Delim('{'):
			key = "object_in_array"
			arr.AddChild(p.parseObject(key))
		case json.Delim('['):
			key = "array_in_array"
			arr.AddChild(p.parseArray(key))
		case json.Delim(']'):
			return arr
		default:
			key = "primitive_in_array"
			arr.AddChild(p.parsePrimitive(key, t))
		}
	}
	return arr
}

func (p *Parser) parsePrimitive(key string, value json.Token) *ds.JSONPrimitive {
	prim := &ds.JSONPrimitive{Key: key}
	switch value.(type) {
	case float64:
		if value == math.Trunc(value.(float64)) {
			prim.Ptype = ds.Int
		} else {
			prim.Ptype = ds.Float
		}
	case string:
		prim.Ptype = ds.String
	case bool:
		prim.Ptype = ds.Bool
	default:
		fmt.Printf("Could not determine datatype of field with key %v and value %v\n", key, value)
		prim.Ptype = ds.String
	}
	return prim
}
