package ds

type JSONElement interface {
	String() string
	Datatype() string
}

type JSONNode interface {
	JSONElement
	AddChild(c JSONElement)
}

type JSONObject struct {
	JSONElement
	JSONNode
	Root     bool
	Key      string
	Children []JSONElement
}

type JSONArray struct {
	JSONElement
	JSONNode
	Root     bool
	Key      string
	Children []JSONElement
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
