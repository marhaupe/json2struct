package ds

import (
	"errors"
	"fmt"
	"strings"
)

func (arr *JSONArray) AddChild(c JSONElement) {
	if arr.Keys == nil {
		arr.Keys = make(map[string]bool)
	}
	ckey := c.GetKey()
	if arr.Keys[ckey] {
		for _, child := range arr.Children {
			if child.Datatype() == Object && child.GetKey() == ckey {
				err := mergeObjects(c, child)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		arr.Children = append(arr.Children, c)
		arr.Keys[ckey] = true
	}
}

func (arr *JSONArray) GetKey() string {
	return arr.Key
}

func (arr *JSONArray) Datatype() Datatype {
	return Array
}

func (arr *JSONArray) String() string {
	foundChildrenTypes := listChildrenTypes(arr.Children)
	var toString string
	if len(foundChildrenTypes) == 1 {
		dataType := foundChildrenTypes[0]
		switch dataType {
		case String, Int, Bool, Float, Null:
			toString = arr.stringPrimitives(dataType)
		case Object:
			toString = arr.stringObjects()
		case Array:
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
	if arr.Root {
		return fmt.Sprintf("type JSONToStruct []struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s []struct{\n%s} `json:\"%s\"`\n", strings.Title(arr.Key), b.String(), arr.Key)
}

func (arr *JSONArray) stringArrays() string {
	if arr.Root {
		return "type JSONToStruct [][]interface{}"
	}
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func (arr *JSONArray) stringPrimitives(dataType Datatype) string {
	typeAsString := arr.Children[0].(*JSONPrimitive).TypeAsString()
	if arr.Root {
		return fmt.Sprintf("type JSONToStruct []%s", typeAsString)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(arr.Key), typeAsString, arr.Key)
}

func (arr *JSONArray) stringMultipleTypes() string {
	if arr.Root {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func listChildrenTypes(c []JSONElement) []Datatype {
	foundChildrenTypes := make(map[Datatype]bool)
	var foundChildren []Datatype
	for _, entry := range c {
		foundChildrenTypes[entry.Datatype()] = true
	}
	for k := range foundChildrenTypes {
		foundChildren = append(foundChildren, k)
	}
	return foundChildren
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
