package ds

import (
	"errors"
	"fmt"
	"strings"
)

func (arr *JSONArray) SetParent(p JSONNode) {
	arr.Parent = p
}

func (arr *JSONArray) GetParent() JSONNode {
	return arr.Parent
}

func (arr *JSONArray) GetKey() string {
	return arr.Key
}

func (arr *JSONArray) GetDatatype() Datatype {
	return Array
}

func (arr *JSONArray) AddChild(c JSONElement) {
	if arr.Types == nil {
		arr.Types = make(map[Datatype]bool)
	}
	ctype := c.GetDatatype()
	if arr.Types[ctype] {
		for _, child := range arr.Children {
			if child.GetDatatype() == Object && child.GetDatatype() == ctype {
				err := mergeObjects(c, child)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		arr.Children = append(arr.Children, c)
		c.SetParent(arr)
		arr.Types[ctype] = true
	}
}

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
		return fmt.Sprintf("type JSONToStruct []struct{\n%s}\n", b.String())
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

func countChildrenTypes(c []JSONElement) int {
	foundChildrenTypes := make(map[Datatype]bool)
	for _, entry := range c {
		foundChildrenTypes[entry.GetDatatype()] = true
	}
	return len(foundChildrenTypes)
}

func mergeObjects(source JSONElement, target JSONElement) error {
	objToAdd, ok := source.(*JSONObject)
	if !ok {
		return errors.New("Error parsing JSONElement to *JSONObject")
	}
	childObj, ok := target.(*JSONObject)
	for _, child := range objToAdd.Children {
		childObj.AddChild(child)
	}
	return nil
}
