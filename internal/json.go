package internal

import (
	"fmt"
	"strings"
)

type PrimitiveType int

const (
	String PrimitiveType = iota
	Integer
	Bool
)

type JSONElement interface {
	String() string
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
	DataType string
	Children []JSONElement
}

type JSONRoot struct {
	JSONElement
	Children []JSONElement
}

func capitalizeKey(k string) string {
	return strings.Title(k)
}

func (jp *JSONPrimitive) String() string {
	var datatype string
	switch jp.Ptype {
	case String:
		datatype = "string"
	case Integer:
		datatype = "int"
	case Bool:
		datatype = "bool"
	}
	return fmt.Sprintf("%s %s `json:\"%s\"`", capitalizeKey(jp.Key), datatype, jp.Key)
}

func (jp *JSONObject) String() string {
	var b strings.Builder
	for _, entry := range jp.Children {
		fmt.Fprintf(&b, entry.String())
	}
	return fmt.Sprintf("%s struct { %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
}

// Arrays are pretty hard since you can basically throw in any valid datatype, e.g.:
// [
//  {
//    "thissucks": true
//  },
//  {
//    "thisdoesntsuck": {
//         "value": false
//    }
//  }
// ]
// Has to become this:
// type JsonToStruct []struct {
// 	Thissucks      bool `json:"thissucks,omitempty"`
// 	Thisdoesntsuck struct {
// 		Value bool `json:"value"`
// 	} `json:"thisdoesntsuck,omitempty"`
// }
// If you're adding different datatypes, e.g [string, int, object], it's still valid json
// TODO!

func (jp *JSONArray) String() string {
	var b strings.Builder
	for _, entry := range jp.Children {
		fmt.Fprintf(&b, entry.String())
	}
	return fmt.Sprintf("%s []struct { %s } `json:\"%s\"`", capitalizeKey(jp.Key), b.String(), jp.Key)
}

func (jp *JSONRoot) String() string {
	var b strings.Builder
	for _, entry := range jp.Children {
		fmt.Fprintf(&b, entry.String())
	}
	return fmt.Sprintf(`type JsonToStruct struct { %s }`, b.String())
}
