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
		Children: make(map[string][]Node),
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
				object.Children[currentKey] = append(object.Children[currentKey], p.parseString())
			}
		case lex.ItemLeftBrace:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseObject())
		case lex.ItemLeftSqrBrace:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseArray())
		case lex.ItemBool:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseBool())
		case lex.ItemNil:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseNil())
		case lex.ItemFloat:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseFloat())
		case lex.ItemInteger:
			object.Children[currentKey] = append(object.Children[currentKey], p.parseInteger())
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
		Children: make([]Node, 0),
	}

	for p.Item = p.Lexer.NextItem(); p.Item.Typ != lex.ItemRightSqrBrace; p.Item = p.Lexer.NextItem() {
		switch p.Item.Typ {
		case lex.ItemLeftBrace:
			array.Children = append(array.Children, p.parseObject())
		case lex.ItemLeftSqrBrace:
			array.Children = append(array.Children, p.parseArray())
		case lex.ItemNil:
			array.Children = append(array.Children, p.parseNil())
		case lex.ItemBool:
			array.Children = append(array.Children, p.parseBool())
		case lex.ItemString:
			array.Children = append(array.Children, p.parseString())
		case lex.ItemInteger:
			array.Children = append(array.Children, p.parseInteger())
		case lex.ItemFloat:
			array.Children = append(array.Children, p.parseFloat())
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

func (p *Parser) parseFloat() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeFloat,
		value:    p.Item.Value,
	}
}

func (p *Parser) parseInteger() *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: NodeTypeInteger,
		value:    p.Item.Value,
	}
}
