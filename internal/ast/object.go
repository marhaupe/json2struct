package ast

import (
	"fmt"
	"strings"
)

// GetKey returns the key of obj
func (obj *JSONObject) GetKey() string {
	return obj.Key
}

// GetParent returns the parent node of obj
func (obj *JSONObject) GetParent() Node {
	return obj.Parent
}

// SetParent sets p as the parent node of obj
func (obj *JSONObject) SetParent(p Node) {
	obj.Parent = p
}

// GetDatatype returns the datatype (Object, Array, Int, String, Float, Bool or Null)
// of obj
func (obj *JSONObject) GetDatatype() Datatype {
	return Object
}

// AddChild adds c to the children of obj and adds obj as parent of c.
// If an child with the same key exists, the child will not be added,
// with the exception of both children being objects. Only then both
// objects will get merged into each other
func (obj *JSONObject) AddChild(c Element) {
	if obj.Keys == nil {
		obj.Keys = make(map[string]bool)
	}
	ckey := c.GetKey()
	if obj.Keys[ckey] {
		for _, child := range obj.Children {
			if child.GetDatatype() == Object && child.GetKey() == ckey {
				childObj := child.(*JSONObject)
				cObj := c.(*JSONObject)
				mergeObjects(cObj, childObj)
			}
		}
	} else {
		obj.Children = append(obj.Children, c)
		c.SetParent(obj)
		obj.Keys[ckey] = true
	}
}

// String returns Go Code ready for unmarshalling
func (obj *JSONObject) String() string {
	var b strings.Builder
	for _, entry := range obj.Children {
		fmt.Fprintf(&b, entry.String())
	}
	if obj.Parent == nil {
		return fmt.Sprintf("type JSONToStruct struct{\n%s}", b.String())
	}
	return fmt.Sprintf("%s struct{\n%s} `json:\"%s\"`\n", strings.Title(obj.Key), b.String(), obj.Key)
}
