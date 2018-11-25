package internal

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/marhaupe/json-to-struct/internal/ds"
)

type Node struct {
	done bool
	el   ds.JSONNode
	prev *Node
	next *Node
}

type Parser struct {
	rootNode    *Node
	currentNode *Node
}

// Parse parsed JSON Tokens received from chan c
func Parse(c chan json.Token, wg *sync.WaitGroup) {
	p := Parser{}
	p.parse(c, wg)
}

// This might help: https://play.golang.org/p/K0cb7hzc6P6
// I could also recursively call parseObject and parseArray in here,
// that might make some things easier (e.g detecting keys)
func (p *Parser) parse(c chan json.Token, wg *sync.WaitGroup) {
	defer wg.Done()

	// Parsing root delim accordingly
	rootNode := &Node{done: false}
	t := <-c
	switch t {
	case json.Delim('{'):
		rootNode.el = &ds.JSONObject{
			Root: true,
		}
	case json.Delim('['):
		rootNode.el = &ds.JSONArray{
			Root: true,
		}
	}
	p.rootNode = rootNode

	var count int
	for t := range c {
		var key string
		if count%2 == 0 {
			fmt.Printf("Key: %v \n", t)
		} else {
			fmt.Printf("Value: %v \n", t)
		}
		switch t {
		case json.Delim('{'): //todo
		case json.Delim('}'): //todo
		case json.Delim('['): //todo
		case json.Delim(']'): //todo
		}

		count++
	}
}
