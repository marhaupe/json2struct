package ds

import (
	"errors"
	"fmt"
	"strings"
)

func (jp *JSONArray) AddChild(c JSONElement) {
	if jp.Keys == nil {
		jp.Keys = make(map[string]bool)
	}
	ckey := c.GetKey()
	if jp.Keys[ckey] {
		for _, child := range jp.Children {
			if child.Datatype() == "object" && child.GetKey() == ckey {
				err := mergeObjects(c, child)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		jp.Children = append(jp.Children, c)
		jp.Keys[ckey] = true
	}
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

func (jp *JSONArray) GetKey() string {
	return jp.Key
}

func (jp *JSONArray) Datatype() string {
	return "array"
}

func (jp *JSONArray) String() string {
	foundChildrenTypes := listChildrenTypes(jp.Children)
	var toString string
	if len(foundChildrenTypes) == 1 {
		dataType := foundChildrenTypes[0]
		switch dataType {
		case "string", "int", "bool", "float64", "interface{}":
			toString = jp.stringPrimitive(dataType)
		case "object":
			toString = jp.stringObject()
		case "array":
			toString = jp.stringArray()
		default:
			panic("Error stringifying array")
		}
	} else {
		toString = jp.stringMultipleTypes()
	}
	return toString
}

func (jp *JSONArray) stringObject() string {
	var b strings.Builder
	for _, child := range jp.Children {
		childObject, ok := child.(*JSONObject)
		if !ok {
			panic(fmt.Sprintf("Error stringifying object %v", child.GetKey()))
		}
		for _, grandchild := range childObject.Children {
			grandchildString := grandchild.String()
			// grandchildString = appendOmitEmptyToRootElement(grandchildString)
			fmt.Fprintf(&b, grandchildString)
		}
	}
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s []struct{\n%s} `json:\"%s\"`\n", strings.Title(jp.Key), b.String(), jp.Key)
}

func (jp *JSONArray) stringArray() string {
	if jp.Root {
		return "type JSONToStruct [][]interface{}"
	}
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`\n", strings.Title(jp.Key), jp.Key)
}

func (jp *JSONArray) stringPrimitive(dataType string) string {
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []%s", dataType)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(jp.Key), dataType, jp.Key)
}

func (jp *JSONArray) stringMultipleTypes() string {
	if jp.Root {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(jp.Key), jp.Key)
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
