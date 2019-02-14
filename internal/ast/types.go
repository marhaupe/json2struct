// Package ast provides datatypes to build up an abstract syntax tree representing
// a JSON object.
package ast

// Datatype represents every valid JSON datatype
type Datatype int

const (
	//Object is the equivalent to a JSON object
	Object Datatype = iota

	//Array is the equivalent to a JSON array
	Array

	//String is the equivalent to a JSON string
	String

	//Int is the equivalent to a JSON int
	Int

	//Float is the equivalent to a JSON float
	Float

	//Bool is the equivalent to a JSON bool
	Bool

	//Null is the equivalent to a JSON object
	Null
)

// Element represents a single element in the JSON tree.
// See JSONObject, JSONArray and JSONPrimitive for a concrete implementation
type Element interface {
	String() string
	GetDatatype() Datatype
	GetParent() Node
	SetParent(p Node)
	GetKey() string
}

// Node represents a single node in the JSON tree. Since a node is a superset
// of an element, it also implements the Element interface.
// See JSONObject or JSONArray for a concrete implementation
type Node interface {
	Element
	AddChild(c Element)
}

// JSONObject represents and mimics the behaviour of a JSON object and implements
// the Node interface
type JSONObject struct {
	Node
	Key      string
	Children []Element
	Parent   Node
	Keys     map[string]bool
}

// JSONArray represents and mimics the behaviour of a JSON array and implements
// the Node interface
type JSONArray struct {
	Node
	Key      string
	Children []Element
	Parent   Node
	Types    map[Datatype]bool
}

// JSONPrimitive represents and mimics the behaviour of a JSON primitive, such as
// String, Int, Float, Bool or Null.
type JSONPrimitive struct {
	Element
	Datatype Datatype
	Parent   Node
	Key      string
}
