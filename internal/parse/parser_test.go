package parse

import (
	"reflect"
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

func mkPrim(typ NodeType, value string) *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: typ,
		value:    value,
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
				json: `{
					"teststring": "hi",
					"testbool": true,
					"testnumber": 5.4
					}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"teststring": []Node{mkPrim(NodeTypeString, "hi")},
					"testbool":   []Node{mkPrim(NodeTypeBool, "true")},
					"testnumber": []Node{mkPrim(NodeTypeNumber, "5.4")},
				},
			)},
		{
			name: "Array with strings",
			args: args{
				json: `[ "test", "1234", "true" ]`,
			},
			want: mkArrayNode(
				[]Node{
					mkPrim(NodeTypeString, "test"),
					mkPrim(NodeTypeString, "1234"),
					mkPrim(NodeTypeString, "true"),
				},
			)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFromString(tt.args.name, tt.args.json); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromString(): \ngot:\n %#v \nwant:\n %#v", got, tt.want)
			}
		})
	}
}
