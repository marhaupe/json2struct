package ds

type JSONElement interface {
	String() string
	Datatype() string
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

type PrimitiveType int

const (
	String PrimitiveType = iota
	Int
	Bool
)

type JSONPrimitive struct {
	JSONElement
	Ptype PrimitiveType
	Key   string
}
