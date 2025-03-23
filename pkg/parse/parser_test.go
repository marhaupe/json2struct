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

func mkPrim(typ NodeType) *PrimitiveNode {
	return &PrimitiveNode{
		NodeType: typ,
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
					mkPrim(NodeTypeString),
					mkPrim(NodeTypeString),
					mkPrim(NodeTypeString),
				},
			),
		},
		{
			name: "Array with different primitives",
			args: args{
				json: `[ "test", 1234, 1.2, true, null ]`,
			},
			want: mkArrayNode(
				[]Node{
					mkPrim(NodeTypeString),
					mkPrim(NodeTypeInteger),
					mkPrim(NodeTypeFloat),
					mkPrim(NodeTypeBool),
					mkPrim(NodeTypeNil),
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
							"teststring": []Node{mkPrim(NodeTypeString)},
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
							mkPrim(NodeTypeString),
							mkPrim(NodeTypeString),
						},
					),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromString(tt.args.json)
			if err != nil {
				t.Errorf("ParseFromString(): got error %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
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
					"testfloat": 5.4,
					"testint": 5,
					"testnil": null
					}`,
			},
			want: mkObjectNode(
				map[string][]Node{
					"teststring": []Node{mkPrim(NodeTypeString)},
					"testbool":   []Node{mkPrim(NodeTypeBool)},
					"testfloat":  []Node{mkPrim(NodeTypeFloat)},
					"testint":    []Node{mkPrim(NodeTypeInteger)},
					"testnil":    []Node{mkPrim(NodeTypeNil)},
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
								"teststring": []Node{mkPrim(NodeTypeString)},
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
								mkPrim(NodeTypeString),
								mkPrim(NodeTypeString),
							},
						),
					},
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromString(tt.args.json)
			if err != nil {
				t.Errorf("ParseFromString(): got error %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromString(): \ngot:\n %#v \nwant:\n %#v", got, tt.want)
			}
		})
	}
}

func TestParseInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "array with comma before closing square brace",
			json: "[ true, ]",
		},
		{
			name: "object with comma before closing curly brace",
			json: `{ "test": "hi", }`,
		},
	}

	for _, test := range tests {
		if _, err := ParseFromString(test.json); err == nil {
			t.Errorf("TestParseInvalidJSON(): expected error, but received none. input: %v", test.json)
		}
	}
}
