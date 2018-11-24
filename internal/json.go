package internal

import (
	"fmt"
	"regexp"
	"strings"
)

type PrimitiveType int

const (
	String PrimitiveType = iota
	Int
	Bool
)

type JSONElement interface {
	String() string
	Datatype() string
}

type JSONPrimitive struct {
	JSONElement
	Ptype PrimitiveType
	Key   string
}

type JSONObject struct {
	JSONElement
	Root     bool
	Key      string
	Children []JSONElement
}

type JSONArray struct {
	JSONElement
	Root     bool
	Key      string
	Children []JSONElement
}

func capitalizeKey(k string) string {
	return strings.Title(k)
}

func (jp *JSONPrimitive) Datatype() string {
	switch jp.Ptype {
	case String:
		return "string"
	case Int:
		return "int"
	case Bool:
		return "bool"
	default:
		return "string"
	}
}
func (jp *JSONPrimitive) String() string {
	return fmt.Sprintf("%s %s `json:\"%s\"`", capitalizeKey(jp.Key), jp.Datatype(), jp.Key)
}

func (jp *JSONObject) Datatype() string {
	return "object"
}

func (jp *JSONObject) String() string {
	var b strings.Builder
	for _, entry := range jp.Children {
		fmt.Fprintf(&b, entry.String())
	}
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct struct{ %s }`", b.String())
	}
	return fmt.Sprintf("%s struct{ %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
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

func (jp *JSONArray) stringMultipleTypes() string {
	if jp.Root {
		return "type JSONToStruct []interface{}"
	}
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`", capitalizeKey(jp.Key), jp.Key)
}

func (jp *JSONArray) stringObject() string {
	var b strings.Builder
	for _, child := range jp.Children {
		childString := child.String()
		childString = appendOmitEmptyToRootElement(childString)
		fmt.Fprintf(&b, childString)
	}
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []struct{ %s }", b.String())
	}
	return fmt.Sprintf("%s []struct{ %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
}

func appendOmitEmptyToRootElement(s string) string {
	re := regexp.MustCompile("`json:\"(.*)\"`$")
	return re.ReplaceAllString(s, "`json:\"$1,omitempty\"`")
}

func (jp *JSONArray) stringArray() string {
	if jp.Root {
		return "type JSONToStruct [][]interface{}"
	}
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`", capitalizeKey(jp.Key), jp.Key)
}

func (jp *JSONArray) stringPrimitive(dataType string) string {
	if jp.Root {
		return fmt.Sprintf("type JSONToStruct []%s", dataType)
	}
	return fmt.Sprintf("%s []%s `json:\"%s\"`", capitalizeKey(jp.Key), dataType, jp.Key)
}
