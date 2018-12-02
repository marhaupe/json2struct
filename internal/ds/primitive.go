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
	case Float:
		return "float64"
	case Bool:
		return "bool"
	case Null:
		return "interface{}"
	default:
		panic(fmt.Sprintf("Datatype of primitive with key %v could not be determined", jp.Key))
	}
}

func (jp *JSONPrimitive) String() string {
	return fmt.Sprintf("%s %s `json:\"%s\"`\n", capitalizeKey(jp.Key), jp.Datatype(), jp.Key)
}
