package ds

import (
	"fmt"
	"strings"
)

func (jp *JSONArray) AddChild(c JSONElement) {
	if jp.Keys == nil {
		jp.Keys = make(map[string]bool)
	}
	key := c.GetKey()
	if jp.Keys[key] {
		for _, child := range jp.Children {
			if child.GetKey() == key && child.Datatype() == "object" {
				casted, ok := c.(*JSONObject)
				if !ok {
					panic("Error parsing JSONElement to *JSONObject")
				}
				castedExistingChild, ok := child.(*JSONObject)
				for _, child := range casted.Children {
					castedExistingChild.AddChild(child)
				}
			}
		}
	} else {
		jp.Children = append(jp.Children, c)
		jp.Keys[key] = true
	}

	// if jp.Keys == nil {
	// 	jp.Keys = make(map[string]bool)
	// }
	// key := c.GetKey()
	// if jp.Keys[key] {
	// 	switch key {
	// 	case "object_in_array":
	// 		// TODO: Joining attributes of objects
	// 		casted, ok := c.(*JSONObject)
	// 		if !ok {
	// 			panic("Error parsing object_in_array to *JSONObject")
	// 		}
	// 		for _, child := range casted.Children {
	// 			casted.AddChild(child)
	// 		}
	// 	case "array_in_array", "primitive_in_array":
	// 		// Don't allow duplicate keys with these datatypes
	// 	default:
	// 		panic(fmt.Sprintf("Error adding child with key %v", key))
	// 	}
	// } else {
	// 	jp.Children = append(jp.Children, c)
	// 	jp.Keys[key] = true
	// }
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
		case "string", "int", "bool":
			toString = jp.stringPrimitive(dataType)
		case "object":
			toString = jp.stringObject()
		case "array":
			toString = jp.stringArray()
		default:
			toString = "error parsing"
		}
	} else {
		toString = jp.stringMultipleTypes()
	}
	return toString
}

func (jp *JSONArray) stringObject() string {
	var b strings.Builder
	for _, child := range jp.Children {
		childString := child.String()
		childString = appendOmitEmptyToRootElement(childString)
		fmt.Fprintf(&b, childString)
	}
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s []struct{\n%s} `json:\"%s\"`\n", capitalizeKey(jp.Key), b.String(), jp.Key)
}

func (jp *JSONArray) stringArray() string {
	if jp.Root {
		return "type JSONToStruct [][]interface{}"
	}
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`\n", capitalizeKey(jp.Key), jp.Key)
}

func (jp *JSONArray) stringPrimitive(dataType string) string {
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []%s", dataType)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`\n", capitalizeKey(jp.Key), dataType, jp.Key)
}

func (jp *JSONArray) stringMultipleTypes() string {
	if jp.Root {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", capitalizeKey(jp.Key), jp.Key)
}
