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
		panic("expected either a curly or a square brace but got something else")
	}
}

func (p *Parser) parseObject() *ObjectNode {
	object := &ObjectNode{NodeType: NodeTypeObject}
	var currentKey string

	for p.Item = p.Lexer.NextItem(); p.Item.Typ != lex.ItemRightBrace; p.Item = p.Lexer.NextItem() {

		if object.children == nil {
			object.children = make(map[string][]Node)
		}

		switch p.Item.Typ {
		case lex.ItemString:

			// A string can indicate that the current lexem is either a key or a value.
			// It's a key if the previous lexem is a comma.
			// It's a value if the previous lexem is a colon.
			if p.LastItem.Typ == lex.ItemComma {
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
			continue
		case lex.ItemComma:
			continue
		default:
			panic("expected a value of type object, array, null, bool, string or number but got something else")
		}
		p.LastItem = p.Item
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("a closing curly brace mustn't follow a comma")
	}

	return object
}

func (p *Parser) parseArray() *ArrayNode {
	array := &ArrayNode{NodeType: NodeTypeArray}

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
		default:
			panic("expected a value of type object, array, null, bool, string or number but got something else")
		}
		p.LastItem = p.Item
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("a closing square brace mustn't follow a comma")
	}

	return array
}

func (p *Parser) parseBool() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeBool,
	}
}

func (p *Parser) parseString() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeString,
	}
}

func (p *Parser) parseNil() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeNil,
	}
}

func (p *Parser) parseNumber() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeNumber,
	}
}
