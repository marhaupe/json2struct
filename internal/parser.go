package internal

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"
)

// type Node struct {
// 	done bool
// 	el   ds.JSONNode
// 	prev *Node
// 	next *Node
// }

type Parser struct {
	// rootNode    *Node
	// currentNode *Node
	rootEl ds.JSONNode
	c      chan json.Token
	wg     *sync.WaitGroup
}

// Parse parsed JSON Tokens received from chan c
func Parse(n chan ds.JSONNode, c chan json.Token, wg *sync.WaitGroup) {
	p := Parser{c: c, wg: wg}
	p.parse()
	n <- p.rootEl
}

// This might help: https://play.golang.org/p/K0cb7hzc6P6
// I could also recursively call parseObject and parseArray in here,
// that might make some things easier (e.g detecting keys)
func (p *Parser) parse() {
	defer p.wg.Done()

	// Consuming first Token, { or [
	t := <-p.c

	// Parsing root delim accordingly
	switch t {
	case json.Delim('{'):
		p.rootEl = &ds.JSONObject{
			Root: true,
		}
	case json.Delim('['):
		p.rootEl = &ds.JSONArray{
			Root: true,
		}
	}

	var key string
	for t := range p.c {
		if key == "" {
			fmt.Printf("Key | Type %T | Value %v\n", t, t)
			key = fmt.Sprint(t)
		} else {
			fmt.Printf("Value | Type %T | Value %v\n", t, t)
			switch t {
			case json.Delim('{'):
				p.rootEl.AddChild(p.parseObject(key))
			case json.Delim('['):
				p.rootEl.AddChild(p.parseArray(key))
			case json.Delim(']'):
				break
			case json.Delim('}'):
				break
			default:
				p.rootEl.AddChild(p.parsePrimitive(key, t))
			}
			key = ""
		}
	}
}

func (p *Parser) parseObject(key string) *ds.JSONObject {
	obj := &ds.JSONObject{}
	return obj
}

func (p *Parser) parseArray(key string) *ds.JSONArray {
	arr := &ds.JSONArray{}
	return arr
}

func (p *Parser) parsePrimitive(key string, value json.Token) *ds.JSONPrimitive {
	prim := &ds.JSONPrimitive{Key: key}
	switch value.(type) {
	case float64:
		prim.Ptype = ds.Int
	case string:
		prim.Ptype = ds.String
	case bool:
		prim.Ptype = ds.Bool
	default:
		fmt.Printf("Could not determine datatype of field with key %v\n", key)
		prim.Ptype = ds.String
	}
	return prim
}
