package ast

import (
	"fmt"
	"strings"
)

// GetKey returns the key of jp
func (jp *JSONPrimitive) GetKey() string {
	return jp.Key
}

// SetParent sets p as the parent node of jp
func (jp *JSONPrimitive) SetParent(p Node) {
	jp.Parent = p
}

// GetParent returns the parent node of jp
func (jp *JSONPrimitive) GetParent() Node {
	return jp.Parent
}

// GetDatatype returns the datatype (Object, Array, Int, String, Float, Bool or Null)
// of jp
func (jp *JSONPrimitive) GetDatatype() Datatype {
	return jp.Datatype
}

// TypeAsString returns the string representation of each possible datatype for jp
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

// String returns Go Code ready for unmarshalling
func (jp *JSONPrimitive) String() string {
	return fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(jp.Key), jp.TypeAsString(), jp.Key)
}
