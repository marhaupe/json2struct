package parse

import (
	"reflect"
	"testing"
)

func mkObjectNode(children map[string][]Node) *ObjectNode {
	return &ObjectNode{
		Children: children,
		NodeType: NodeTypeObject,
	}
}

func mkArrayNode(children []Node) *ArrayNode {
	return &ArrayNode{
		Children: children,
		NodeType: NodeTypeArray,
	}
}

func mkPrim(typ NodeType, value string) *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: typ,
		value:    value,
	}
}

func TestParseArrayFromString(t *testing.T) {
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
			name: "Empty array",
			args: args{
				json: `[]`,
			},
			want: mkArrayNode(make([]Node, 0)),
		},
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
			),
		},
		{
			name: "Array with different primitives",
			args: args{
				json: `[ "test", 1234, true, null ]`,
			},
			want: mkArrayNode(
				[]Node{
					mkPrim(NodeTypeString, "test"),
					mkPrim(NodeTypeInteger, "1234"),
					mkPrim(NodeTypeBool, "true"),
					mkPrim(NodeTypeNil, "null"),
				},
			),
		},
		{
			name: "Array with objects",
			args: args{
				json: `[{ "teststring": "hi" }]`,
			},
			want: mkArrayNode(
				[]Node{
					mkObjectNode(
						map[string][]Node{
							"teststring": []Node{mkPrim(NodeTypeString, "hi")},
						},
					),
				},
			),
		},
		{
			name: "Array with arrays",
			args: args{
				json: `[[ "hi", "ho" ]]`,
			},
			want: mkArrayNode(
				[]Node{
					mkArrayNode(
						[]Node{
							mkPrim(NodeTypeString, "hi"),
							mkPrim(NodeTypeString, "ho"),
						},
					),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFromString(tt.args.name, tt.args.json); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromString(): \ngot:\n %#v \nwant:\n %#v", got, tt.want)
			}
		})
	}
}

func TestParseObjectFromString(t *testing.T) {
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
			name: "Object with primitives",
			args: args{
				json: `{
					"teststring": "hi",
					"testbool": true,
					"testnumber": 5.4,
					"testnil": null
					}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"teststring": []Node{mkPrim(NodeTypeString, "hi")},
					"testbool":   []Node{mkPrim(NodeTypeBool, "true")},
					"testnumber": []Node{mkPrim(NodeTypeFloat, "5.4")},
					"testnil":    []Node{mkPrim(NodeTypeNil, "null")},
				},
			),
		},
		{
			name: "Object with objects",
			args: args{
				json: `{ "testobject": { "teststring": "hi" }}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"testobject": []Node{
						mkObjectNode(
							map[string][]Node{
								"teststring": []Node{mkPrim(NodeTypeString, "hi")},
							},
						),
					},
				},
			),
		},
		{
			name: "Object with arrays",
			args: args{
				json: `{ "testarray": [ "hi", "ho" ]}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"testarray": []Node{
						mkArrayNode(
							[]Node{
								mkPrim(NodeTypeString, "hi"),
								mkPrim(NodeTypeString, "ho"),
							},
						),
					},
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFromString(tt.args.name, tt.args.json); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromString(): \ngot:\n %#v \nwant:\n %#v", got, tt.want)
			}
		})
	}
}
