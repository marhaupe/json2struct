package json

import (
	"fmt"
	"strings"
)

type JSONObject struct {
	JSONElement
	Root     bool
	Key      string
	Children []JSONElement
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
		return fmt.Sprintf("type JSONToStruct struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s struct{\n%s} `json:\"%s\"`\n", capitalizeKey(jp.Key), b.String(), jp.Key)
}
