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
	Key      string
	Children []JSONElement
}

type JSONArray struct {
	JSONElement
	Key      string
	Children []JSONElement
}

type JSONRoot struct {
	JSONElement
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
	return fmt.Sprintf("%s struct{ %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
}

func (jp *JSONArray) Datatype() string {
	return "array"
}

func (jp *JSONArray) String() string {
	var lastFoundDataType string
	foundChildrenTypes := make(map[string]bool)
	for _, entry := range jp.Children {
		foundChildrenTypes[entry.Datatype()] = true
		lastFoundDataType = entry.Datatype()
	}
	var toString string
	if len(foundChildrenTypes) == 1 {
		dataType := lastFoundDataType
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

func (jp *JSONArray) stringMultipleTypes() string {
	return fmt.Sprintf("%s []interface{} `json:\"%s\"`", capitalizeKey(jp.Key), jp.Key)
}

func (jp *JSONArray) stringObject() string {
	var b strings.Builder
	for _, child := range jp.Children {
		childString := child.String()
		childString = appendOmitEmptyToRootElement(childString)
		fmt.Fprintf(&b, childString)
	}
	return fmt.Sprintf("%s []struct{ %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
}

func appendOmitEmptyToRootElement(s string) string {
	re := regexp.MustCompile("`json:\"(.*)\"`$")
	return re.ReplaceAllString(s, "`json:\"$1,omitempty\"`")
}

func (jp *JSONArray) stringArray() string {
	return fmt.Sprintf("%s [][]interface{} `json:\"%s\"`", capitalizeKey(jp.Key), jp.Key)
}
func (jp *JSONArray) stringPrimitive(dataType string) string {
	return fmt.Sprintf("%s []%s `json:\"%s\"`", capitalizeKey(jp.Key), dataType, jp.Key)
}
func (jp *JSONRoot) Datatype() string {
	return "root"
}

func (jp *JSONRoot) String() string {
	var b strings.Builder
	for _, entry := range jp.Children {
		fmt.Fprintf(&b, entry.String())
	}
	return fmt.Sprintf(`type JsonToStruct struct{ %s }`, b.String())
}
