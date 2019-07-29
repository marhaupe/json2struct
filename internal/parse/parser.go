package parse

import "github.com/marhaupe/json2struct/internal/lex"

type Parser struct {
	Lexer    *lex.Lexer
	Item     lex.Item
	LastItem lex.Item
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
	NodeTypeKey
	NodeTypeString
	NodeTypeBool
	NodeTypeNil
	NodeTypeNumber
)

type ArrayNode struct {
	NodeType
	children []Node
}

type ObjectNode struct {
	NodeType

	// The JSON spec allows different types of values for the same key. Because of that, a simple `map[string]Node` is not enough.
	children map[string][]Node
}

type PrimitiveNode struct {
	NodeType
	value string
}

func ParseFromString(name, json string) Node {
	parser := &Parser{
		Lexer: lex.Lex(name, json),
	}
	return parser.parse()
}

func (p *Parser) parse() Node {
	p.Item = p.Lexer.NextItem()

	switch p.Item.Typ {
	case lex.ItemLeftBrace:
		return p.parseObject()
	case lex.ItemLeftSqrBrace:
		return p.parseArray()
	default:
		panic("error parsing document: unexpected item " + p.Item.Value)
	}
}

func (p *Parser) parseObject() *ObjectNode {
	p.LastItem = p.Item

	object := &ObjectNode{
		NodeType: NodeTypeObject,
		children: make(map[string][]Node),
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
				object.children[currentKey] = append(object.children[currentKey], p.parseString())
			}
		case lex.ItemLeftBrace:
			object.children[currentKey] = append(object.children[currentKey], p.parseObject())
		case lex.ItemLeftSqrBrace:
			object.children[currentKey] = append(object.children[currentKey], p.parseArray())
		case lex.ItemBool:
			object.children[currentKey] = append(object.children[currentKey], p.parseBool())
		case lex.ItemNil:
			object.children[currentKey] = append(object.children[currentKey], p.parseNil())
		case lex.ItemNumber:
			object.children[currentKey] = append(object.children[currentKey], p.parseNumber())
		case lex.ItemColon:
			break
		case lex.ItemComma:
			break
		default:
			panic("error parsing object: unexpected item " + p.Item.Value)
		}
		p.LastItem = p.Item
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("error parsing object: a closing curly brace mustn't follow a comma")
	}

	return object
}

func (p *Parser) parseArray() *ArrayNode {
	p.LastItem = p.Item
	array := &ArrayNode{
		NodeType: NodeTypeArray,
		children: make([]Node, 0),
	}

	for p.Item = p.Lexer.NextItem(); p.Item.Typ != lex.ItemRightSqrBrace; p.Item = p.Lexer.NextItem() {
		switch p.Item.Typ {
		case lex.ItemLeftBrace:
			array.children = append(array.children, p.parseObject())
		case lex.ItemLeftSqrBrace:
			array.children = append(array.children, p.parseArray())
		case lex.ItemNil:
			array.children = append(array.children, p.parseNil())
		case lex.ItemBool:
			array.children = append(array.children, p.parseBool())
		case lex.ItemString:
			array.children = append(array.children, p.parseString())
		case lex.ItemNumber:
			array.children = append(array.children, p.parseNumber())
		case lex.ItemComma:
			break
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

func (p *Parser) parseBool() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeBool,
		value:    p.Item.Value,
	}
}

func (p *Parser) parseString() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeString,
		value:    p.Item.Value,
	}
}

func (p *Parser) parseNil() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeNil,
		value:    p.Item.Value,
	}
}

func (p *Parser) parseNumber() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeNumber,
		value:    p.Item.Value,
	}
}
