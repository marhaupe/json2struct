package ds

import (
	"fmt"
	"strings"
)

func (jp *JSONObject) GetKey() string {
	return jp.Key
}

func (jp *JSONObject) Datatype() string {
	return "object"
}

func (jp *JSONObject) AddChild(c JSONElement) {
	if jp.Keys == nil {
		jp.Keys = make(map[string]bool)
	}
	key := c.GetKey()
	if jp.Keys[key] {
		for _, child := range jp.Children {
			if child.GetKey() == key && child.Datatype() == "object" {
				castedToAdd, ok := c.(*JSONObject)
				if !ok {
					panic("Error parsing JSONElement to *JSONObject")
				}
				castedExistingChild, ok := child.(*JSONObject)
				if !ok {
					panic("Error parsing JSONElement to *JSONObject")
				}
				for _, child := range castedToAdd.Children {
					castedExistingChild.AddChild(child)
				}
			}
		}
	} else {
		jp.Children = append(jp.Children, c)
		jp.Keys[key] = true
	}
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
