package ds

import (
	"fmt"
	"strings"
)

func (obj *JSONObject) GetKey() string {
	return obj.Key
}

func (obj *JSONObject) Datatype() string {
	return "object"
}

func (obj *JSONObject) AddChild(c JSONElement) {
	if obj.Keys == nil {
		obj.Keys = make(map[string]bool)
	}
	ckey := c.GetKey()
	if obj.Keys[ckey] {
		for _, child := range obj.Children {
			if child.Datatype() == "object" && child.GetKey() == ckey {
				err := mergeObjects(c, child)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		obj.Children = append(obj.Children, c)
		obj.Keys[ckey] = true
	}
}

func (obj *JSONObject) String() string {
	var b strings.Builder
	for _, entry := range obj.Children {
		fmt.Fprintf(&b, entry.String())
	}
	if obj.Root {
		return fmt.Sprintf("type JSONToStruct struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s struct{\n%s} `json:\"%s\"`\n", strings.Title(obj.Key), b.String(), obj.Key)
}
