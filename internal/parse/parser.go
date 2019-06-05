package parse

type Node interface {
	Type() NodeType
	String() string
	Copy() Node
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeIdent NodeType = iota
	NodeString
	NodeNumber
	NodeCall
	NodeVector
)

func ParseFromString(s string) Node {
	return nil
}
