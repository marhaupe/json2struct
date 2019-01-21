package tree

import (
	"fmt"
	"strings"
)

func (obj *JSONObject) GetKey() string {
	return obj.Key
}

func (obj *JSONObject) GetParent() JSONNode {
	return obj.Parent
}

func (obj *JSONObject) SetParent(p JSONNode) {
	obj.Parent = p
}

func (obj *JSONObject) GetDatatype() Datatype {
	return Object
}

func (obj *JSONObject) AddChild(c JSONElement) {
	if obj.Keys == nil {
		obj.Keys = make(map[string]bool)
	}
	ckey := c.GetKey()
	if obj.Keys[ckey] {
		for _, child := range obj.Children {
			if child.GetDatatype() == Object && child.GetKey() == ckey {
				childObj := child.(*JSONObject)
				cObj := c.(*JSONObject)
				mergeObjects(cObj, childObj)
			}
		}
	} else {
		obj.Children = append(obj.Children, c)
		c.SetParent(obj)
		obj.Keys[ckey] = true
	}
}

func (obj *JSONObject) String() string {
	var b strings.Builder
	for _, entry := range obj.Children {
		fmt.Fprintf(&b, entry.String())
	}
	if obj.Parent == nil {
		return fmt.Sprintf("type JSONToStruct struct{\n%s}\n", b.String())
	}
	return fmt.Sprintf("%s struct{\n%s} `json:\"%s\"`\n", strings.Title(obj.Key), b.String(), obj.Key)
}
