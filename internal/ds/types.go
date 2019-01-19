package ds

type Datatype int

const (
	Object Datatype = iota
	Array
	String
	Int
	Float
	Bool
	Null
)

type JSONElement interface {
	String() string
	GetDatatype() Datatype
	GetKey() string
}

type JSONNode interface {
	AddChild(c JSONElement)
}

type JSONObject struct {
	Root     bool
	Key      string
	Children []JSONElement
	Keys     map[string]bool
}

type JSONArray struct {
	Root     bool
	Key      string
	Children []JSONElement
	Keys     map[string]bool
}

type JSONPrimitive struct {
	Datatype Datatype
	Key      string
}
