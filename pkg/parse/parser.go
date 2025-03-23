package parse

import (
	"errors"
	"fmt"

	"github.com/marhaupe/json2struct/pkg/lex"
)

type Parser struct {
	Lexer    *lex.Lexer
	Item     *lex.Item
	LastItem *lex.Item
}

type Node interface {
	Type() NodeType
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeTypeArray NodeType = iota
	NodeTypeObject
	NodeTypeString
	NodeTypeBool
	NodeTypeNil
	NodeTypeFloat
	NodeTypeInteger
)

type ArrayNode struct {
	NodeType
	Children []Node
}

type ObjectNode struct {
	NodeType
	// The JSON spec allows different types of values for the same key. Because of that, a simple `map[string]Node` is not enough.
	Children map[string][]Node
}

type PrimitiveNode struct {
	NodeType
}

// Singleton primitive nodes to avoid allocations
var (
	stringNode  = &PrimitiveNode{NodeType: NodeTypeString}
	boolNode    = &PrimitiveNode{NodeType: NodeTypeBool}
	nilNode     = &PrimitiveNode{NodeType: NodeTypeNil}
	floatNode   = &PrimitiveNode{NodeType: NodeTypeFloat}
	integerNode = &PrimitiveNode{NodeType: NodeTypeInteger}
)

func ParseFromString(j string) (Node, error) {
	parser := &Parser{
		Lexer: lex.Lex(j),
	}
	return parser.parse()
}

func (p *Parser) parse() (node Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = errors.New(fmt.Sprint(r))
		}
	}()

	p.Item = p.Lexer.NextItem()

	switch p.Item.Typ {
	case lex.ItemLeftBrace:
		return p.parseObject(), nil
	case lex.ItemLeftSqrBrace:
		return p.parseArray(), nil
	case lex.ItemError:
		panic(fmt.Sprintf("received error from lexer at pos %v: %v", p.Item.Pos, p.Item.Value))
	default:
		panic(fmt.Sprintf("error determining root json type. unexpected item %v", p.Item.Value))
	}
}

func (p *Parser) parseObject() *ObjectNode {
	p.LastItem = p.Item

	object := &ObjectNode{
		NodeType: NodeTypeObject,
		Children: make(map[string][]Node, 8),
	}

	var currentKey string
	for p.Item = p.Lexer.NextItem(); p.Item.Typ != lex.ItemRightBrace; p.Item = p.Lexer.NextItem() {

		switch p.Item.Typ {
		case lex.ItemString:

			// A string can indicate that the current lexem is either a key or a value.
			// It's a key if the previous lexem is a comma or an opening brace.
			// It's a value if the previous lexem is a colon.
			if p.LastItem.Typ == lex.ItemComma || p.LastItem.Typ == lex.ItemLeftBrace {
				currentKey = p.Item.Value
			} else {
				object.Children[currentKey] = append(object.Children[currentKey], stringNode)
			}
		case lex.ItemLeftBrace:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseObject())
		case lex.ItemLeftSqrBrace:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseArray())
		case lex.ItemBool:
			object.Children[currentKey] = append(object.Children[currentKey], boolNode)
		case lex.ItemNil:
			object.Children[currentKey] = append(object.Children[currentKey], nilNode)
		case lex.ItemFloat:
			object.Children[currentKey] = append(object.Children[currentKey], floatNode)
		case lex.ItemInteger:
			object.Children[currentKey] = append(object.Children[currentKey], integerNode)
		case lex.ItemColon:
			break
		case lex.ItemComma:
			break
		case lex.ItemError:
			panic(fmt.Sprintf("received error from lexer. pos: %v, msg: %v", p.Item.Pos, p.Item.Value))
		default:
			panic(fmt.Sprintf("error parsing object. unexpected item %v", p.Item.Value))
		}
		p.LastItem = p.Item
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("error parsing object. a closing curly brace mustn't follow a comma")
	}

	return object
}

func (p *Parser) parseArray() *ArrayNode {
	p.LastItem = p.Item
	array := &ArrayNode{
		NodeType: NodeTypeArray,
		Children: make([]Node, 0, 8),
	}

	for p.Item = p.Lexer.NextItem(); p.Item.Typ != lex.ItemRightSqrBrace; p.Item = p.Lexer.NextItem() {
		switch p.Item.Typ {
		case lex.ItemLeftBrace:
			array.Children = append(array.Children, p.parseObject())
		case lex.ItemLeftSqrBrace:
			array.Children = append(array.Children, p.parseArray())
		case lex.ItemNil:
			array.Children = append(array.Children, nilNode)
		case lex.ItemBool:
			array.Children = append(array.Children, boolNode)
		case lex.ItemString:
			array.Children = append(array.Children, stringNode)
		case lex.ItemInteger:
			array.Children = append(array.Children, integerNode)
		case lex.ItemFloat:
			array.Children = append(array.Children, floatNode)
		case lex.ItemComma:
			break
		case lex.ItemError:
			panic(fmt.Sprintf("received error from lexer. pos: %v, msg: %v", p.Item.Pos, p.Item.Value))
		default:
			panic("error parsing array: unexpected item " + p.Item.Value)
		}
		p.LastItem = p.Item
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("error parsing array: a closing square brace mustn't follow a comma")
	}

	return array
}
