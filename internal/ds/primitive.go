package ds

import "fmt"

func (jp *JSONPrimitive) GetKey() string {
	return jp.Key
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
	return fmt.Sprintf("%s %s `json:\"%s\"`\n", capitalizeKey(jp.Key), jp.Datatype(), jp.Key)
}
