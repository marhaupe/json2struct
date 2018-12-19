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
			if child.Datatype() == "object" && child.GetKey() == ckey {
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

func (arr *JSONArray) Datatype() string {
	return "array"
}

func (arr *JSONArray) String() string {
	foundChildrenTypes := listChildrenTypes(arr.Children)
	var toString string
	if len(foundChildrenTypes) == 1 {
		dataType := foundChildrenTypes[0]
		switch dataType {
		case "string", "int", "bool", "float64", "interface{}":
			toString = arr.stringPrimitive(dataType)
		case "object":
			toString = arr.stringObject()
		case "array":
			toString = arr.stringArray()
		default:
			panic("Error stringifying array")
		}
	} else {
		toString = arr.stringMultipleTypes()
	}
	return toString
}

func (arr *JSONArray) stringObject() string {
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

func (arr *JSONArray) stringArray() string {
	if arr.Root {
		return "type JSONToStruct [][]interface{}"
	}
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func (arr *JSONArray) stringPrimitive(dataType string) string {
	if arr.Root {
		return fmt.Sprintf("type JSONToStruct []%s", dataType)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(arr.Key), dataType, arr.Key)
}

func (arr *JSONArray) stringMultipleTypes() string {
	if arr.Root {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(arr.Key), arr.Key)
}

func listChildrenTypes(c []JSONElement) []string {
	foundChildrenTypes := make(map[string]bool)
	var foundChildren []string
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
