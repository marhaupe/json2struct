package ast

import (
	"fmt"
	"strings"
)

// SetParent sets p as the parent node of arr
func (arr *JSONArray) SetParent(p Node) {
	arr.Parent = p
}

// GetParent returns the parent node of arr
func (arr *JSONArray) GetParent() Node {
	return arr.Parent
}

// GetKey returns the key of arr
func (arr *JSONArray) GetKey() string {
	return arr.Key
}

// GetDatatype returns the datatype (Object, Array, Int, String, Float, Bool or Null)
// of arr
func (arr *JSONArray) GetDatatype() Datatype {
	return Array
}

// AddChild adds c to the children of arr and adds arr as parent of c.
// If an child with the same key exists, the child will not be added,
// with the exception of both children being objects. Only then both
// objects will get merged into each other
func (arr *JSONArray) AddChild(c Element) {
	if arr.Types == nil {
		arr.Types = make(map[Datatype]bool)
	}
	ctype := c.GetDatatype()
	if arr.Types[ctype] && ctype == Object {
		for _, child := range arr.Children {
			if child.GetDatatype() == Object {
				childObj := child.(*JSONObject)
				cObj := c.(*JSONObject)
				mergeObjects(cObj, childObj)
			}
		}
	} else {
		arr.Children = append(arr.Children, c)
		c.SetParent(arr)
		arr.Types[ctype] = true
	}
}

// String returns Go Code ready for unmarshalling
func (arr *JSONArray) String() string {
	childrenTypeCount := countChildrenTypes(arr.Children)
	var toString string
	if childrenTypeCount == 1 {
		firstChild := arr.Children[0]
		switch firstChild.(type) {
		case *JSONPrimitive:
			toString = arr.stringPrimitives(firstChild.(*JSONPrimitive).TypeAsString())
		case *JSONObject:
			toString = arr.stringObjects()
		case *JSONArray:
			toString = arr.stringArrays()
		default:
			panic("Error stringifying array")
		}
	} else {
		toString = arr.stringMultipleTypes()
	}
	return toString
}

func (arr *JSONArray) stringObjects() string {
	var b strings.Builder
	for _, child := range arr.Children {
		childObject, ok := child.(*JSONObject)
		if !ok {
			panic(fmt.Sprintf("Error stringifying object %v", child.GetKey()))
		}
		for _, grandchild := range childObject.Children {
			grandchildString := grandchild.String()
			fmt.Fprintf(&b, grandchildString)
		}
	}
	if arr.Parent == nil {
		return fmt.Sprintf("type JSONToStruct []struct{\n%s}", b.String())
	}
	return fmt.Sprintf("%s []struct{\n%s} `json:\"%s\"`\n", strings.Title(arr.Key), b.String(), arr.Key)
}

func (arr *JSONArray) stringArrays() string {
	if arr.Parent == nil {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func (arr *JSONArray) stringPrimitives(dataType string) string {
	if arr.Parent == nil {
		return fmt.Sprintf("type JSONToStruct []%s", dataType)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(arr.Key), dataType, arr.Key)
}

func (arr *JSONArray) stringMultipleTypes() string {
	if arr.Parent == nil {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func countChildrenTypes(c []Element) int {
	foundChildrenTypes := make(map[Datatype]bool)
	for _, entry := range c {
		foundChildrenTypes[entry.GetDatatype()] = true
	}
	return len(foundChildrenTypes)
}

func mergeObjects(source *JSONObject, target *JSONObject) {
	for _, child := range source.Children {
		target.AddChild(child)
	}
}
