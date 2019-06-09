package parse

import "github.com/marhaupe/json2struct/internal/lex"

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
	NodeTypeKey
	NodeTypeString
	NodeTypeBool
	NodeTypeNil
	NodeTypeFloat
)

type ArrayNode struct {
	NodeType
	children []Node
}

type ObjectNode struct {
	NodeType
	children map[string]Node
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
	lexem := p.Lexer.NextItem()
	p.LastItem = &lexem

	switch lexem.Typ {
	case lex.ItemLeftBrace:
		return p.parseObject()
	case lex.ItemLeftSqrBrace:
		return p.parseArray()
	default:
		panic("expected either a curly or a square brace but got something else")
	}
}

func (p *Parser) parseObject() Node {
	var object *ObjectNode
	var currentKey string

	for lexem := p.Lexer.NextItem(); lexem.Typ != lex.ItemRightBrace; lexem = p.Lexer.NextItem() {
		switch lexem.Typ {
		case lex.ItemString:

			// a string can indicate that the current lexem is either a key or a value. it's a key if the previous lexem is a  comma. It's a value if the
			// previous lexem is a colon.
			if p.LastItem.Typ == lex.ItemComma {
				currentKey = lexem.Value
			} else {
				object.children[currentKey] = p.parsePrimitive()
			}
		case lex.ItemLeftBrace:
			object.children[currentKey] = p.parseObject()
		case lex.ItemLeftSqrBrace:
			object.children[currentKey] = p.parseArray()
		case lex.ItemBool:
			fallthrough
		case lex.ItemNil:
			fallthrough
		case lex.ItemNumber:
			object.children[currentKey] = p.parsePrimitive()
		default:
			panic("expected a value of type object, array, null, bool, string or number but got something else")
		}
		p.LastItem = &lexem
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("a closing curly brace mustn't follow a comma")
	}

	return object
}

func (p *Parser) parseArray() *ArrayNode {
	var array *ArrayNode

	for lexem := p.Lexer.NextItem(); lexem.Typ != lex.ItemRightSqrBrace; lexem = p.Lexer.NextItem() {
		switch lexem.Typ {
		case lex.ItemLeftBrace:
			array.children = append(array.children, p.parseObject())
		case lex.ItemLeftSqrBrace:
			array.children = append(array.children, p.parseArray())
		case lex.ItemNil:
			fallthrough
		case lex.ItemBool:
			fallthrough
		case lex.ItemString:
			fallthrough
		case lex.ItemNumber:
			array.children = append(array.children, p.parsePrimitive())
		default:
			panic("expected a value of type object, array, null, bool, string or number but got something else")
		}
		p.LastItem = &lexem
	}

	if p.LastItem.Typ == lex.ItemComma {
		panic("a closing square brace mustn't follow a comma")
	}

	return array
}

func (p *Parser) parsePrimitive() *PrimitiveNode {
	return nil
}
