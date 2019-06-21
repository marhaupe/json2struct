package parse

import (
	"testing"
)

func mkObjectNode(children map[string][]Node) *ObjectNode {
	return &ObjectNode{
		children: children,
		NodeType: NodeTypeObject,
	}
}

func mkArrayNode(children []Node) *ArrayNode {
	return &ArrayNode{
		children: children,
		NodeType: NodeTypeArray,
	}
}

func mkPrim(typ NodeType) *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: typ,
	}
}

func equal(t1, t2 Node) bool {
	if t1.Type() != t2.Type() {
		return false
	}

	switch t1.Type() {
	case NodeTypeArray:
		castedt1 := t1.(*ArrayNode)
		castedt2 := t2.(*ArrayNode)

		if len(castedt1.children) != len(castedt2.children) {
			return false
		}

		for i := 0; i < len(castedt1.children); i++ {

			if !equal(castedt1.children[i], castedt2.children[i]) {
				return false
			}
		}
		return true

	case NodeTypeObject:
		castedt1 := t1.(*ObjectNode)
		castedt2 := t2.(*ObjectNode)

		if len(castedt1.children) != len(castedt2.children) {
			return false
		}

		for key := range castedt1.children {
			if len(castedt1.children[key]) != len(castedt2.children[key]) {
				return false
			}

			for index := range castedt1.children[key] {
				if !equal(castedt2.children[key][index], castedt2.children[key][index]) {
					return false
				}
			}
		}
		return true

	default:
		return true
	}
}

func TestParseFromString(t *testing.T) {
	type args struct {
		name string
		json string
	}
	tests := []struct {
		name string
		args args
		want Node
	}{
		{
			name: "Empty object",
			args: args{
				json: `{}`,
			},
			want: mkObjectNode(make(map[string][]Node)),
		},
		{
			name: "Empty array",
			args: args{
				json: `[]`,
			},
			want: mkArrayNode(make([]Node, 0)),
		},
		{
			name: "Object with primitives",
			args: args{
				// TODO: This actually does not work the way I want it to work.
				// Ideally, the generated tree contains the value aswell.
				json: `{
					"teststring": "hi",
					"testbool": true,
					"testnumber": 5.4
					}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"teststring": []Node{PrimitiveNode{NodeType: NodeTypeString}},
					"testbool":   []Node{PrimitiveNode{NodeType: NodeTypeBool}},
					"testnumber": []Node{PrimitiveNode{NodeType: NodeTypeNumber}},
				},
			)},
		{
			name: "Array with strings",
			args: args{
				json: `[ "test", "1234", "true" ]`,
			},
			want: mkArrayNode(
				[]Node{
					PrimitiveNode{NodeType: NodeTypeString},
					PrimitiveNode{NodeType: NodeTypeString},
					PrimitiveNode{NodeType: NodeTypeString},
				},
			)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFromString(tt.args.name, tt.args.json); !equal(got, tt.want) {
				t.Errorf("ParseFromString(): \ngot:\n %#v \nwant:\n %#v", got, tt.want)
			}
		})
	}
}
