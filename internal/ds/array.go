package ds

import (
	"fmt"
	"strings"
)

type JSONArray struct {
	JSONElement
	JSONNode
	Root     bool
	Key      string
	Children []JSONElement
}

func (jp *JSONArray) AddChild(c JSONElement) {
	jp.Children = append(jp.Children, c)
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
