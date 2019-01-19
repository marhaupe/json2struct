package ds

import (
	"fmt"
	"strings"
)

func (jp *JSONPrimitive) GetKey() string {
	return jp.Key
}

func (jp *JSONPrimitive) SetParent(p JSONNode) {
	jp.Parent = p
}

func (jp *JSONPrimitive) GetParent() JSONNode {
	return jp.Parent
}

func (jp *JSONPrimitive) GetDatatype() Datatype {
	return jp.Datatype
}

func (jp *JSONPrimitive) TypeAsString() string {
	switch jp.Datatype {
	case String:
		return "string"
	case Int:
		return "int"
	case Bool:
		return "bool"
	case Float:
		return "float64"
	case Null:
		return "interface{}"
	default:
		panic("TypeAsString could not detect Type properly")
	}
}

func (jp *JSONPrimitive) String() string {
	return fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(jp.Key), jp.TypeAsString(), jp.Key)
}
