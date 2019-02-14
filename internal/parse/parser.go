// Package parse provides fields to enable the parsing of a JSON, meant to be used in
// conjunction with package lex
package parse

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/marhaupe/json2struct/internal/ast"
	"github.com/marhaupe/json2struct/internal/lex"
)

// Parser contains fields to enable parsing of JSON tokens to an ast
type Parser struct {
	rootEl ast.Node
	lexRes chan lex.Result
	wg     *sync.WaitGroup
}

// Result is the result of a the parsing of JSON. It either contains Node, representing
// the root JSON node and thus the root of the ast, or an error
type Result struct {
	Node  ast.Node
	Error error
}

// Parse parses JSON tokens received from lexRes and writes either the resulting error
// or the parsed tree to parseRes
func Parse(parseRes chan Result, lexRes chan lex.Result, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			parseRes <- Result{Node: nil, Error: err}
			close(parseRes)
		}
	}()

	p := Parser{lexRes: lexRes, wg: wg}
	p.parse()
	parseRes <- Result{Node: p.rootEl, Error: nil}
	close(parseRes)
}

func (p *Parser) parse() {
	defer p.wg.Done()

	// Consuming first Token, { or [
	r := <-p.lexRes
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

func (p *Parser) parseObject(objKey string) *ast.JSONObject {
	obj := &ast.JSONObject{Key: objKey}
	var key string
	for r := range p.lexRes {
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
			done := p.buildUpElement(obj, key, r.Token)
			if done {
				return obj
			}
			key = ""
		}
	}
	return obj
}

func (p *Parser) parseArray(arrKey string) *ast.JSONArray {
	arr := &ast.JSONArray{Key: arrKey}

	for t := range p.lexRes {
		if t.Error != nil {
			panic(t.Error)
		}
		key := ""
		done := p.buildUpElement(arr, key, t.Token)
		if done {
			return arr
		}
	}
	return arr
}

func (p *Parser) buildUpElement(node ast.Node, key string, t json.Token) (done bool) {
	switch t {
	case json.Delim('}'), json.Delim(']'):
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

func (p *Parser) parsePrimitive(key string, value json.Token) *ast.JSONPrimitive {
	prim := &ast.JSONPrimitive{Key: key}
	prim.Datatype = detectPrimitiveType(value)
	return prim
}

func detectPrimitiveType(t json.Token) ast.Datatype {
	switch t.(type) {
	case float64:
		if t == math.Trunc(t.(float64)) {
			return ast.Int
		}
		return ast.Float
	case string:
		return ast.String
	case bool:
		return ast.Bool
	case nil:
		return ast.Null
	default:
		fmt.Printf("Could not determine datatype of field with value %v\n", t)
		return ast.String
	}
}
