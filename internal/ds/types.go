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
	GetParent() JSONNode
	SetParent(p JSONNode)
	GetKey() string
}

type JSONNode interface {
	AddChild(c JSONElement)
}

type JSONObject struct {
	Key      string
	Children []JSONElement
	Parent   JSONNode
	Keys     map[string]bool
}

type JSONArray struct {
	Key      string
	Children []JSONElement
	Parent   JSONNode
	Types    map[Datatype]bool
}

type JSONPrimitive struct {
	Datatype Datatype
	Parent   JSONNode
	Key      string
}
