package ds

import (
	"reflect"
	"testing"
)

func TestJSONArray_String(t *testing.T) {
	type fields struct {
		Key      string
		Parent   JSONNode
		Children []JSONElement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Array Test Only Bools",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Bool,
						Key:      "testbool",
					},
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Bool,
						Key:      "testbool2",
					},
				},
			},
			want: "Testarray []bool `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Only Strings",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: String,
						Key:      "teststring1",
					},
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: String,
						Key:      "teststring2",
					},
				},
			},
			want: "Testarray []string `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Only Ints",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Int,
						Key:      "testint1",
					},
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Int,
						Key:      "testint2",
					},
				},
			},
			want: "Testarray []int `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Arrays",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONArray{
						Children: []JSONElement{
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: Bool,
							},
							&JSONPrimitive{
								Datatype: Int,
								Parent:   &JSONArray{},
							},
						}},
					&JSONArray{
						Children: []JSONElement{
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: Bool,
							},
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: Int,
							},
						},
					},
				},
			},
			want: "Testarray []interface{} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Objects",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONObject{
						Key: "Intstringobj",
						Children: []JSONElement{
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: Int,
								Key:      "testint",
							},
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: String,
								Key:      "teststring",
							},
						},
					},
					&JSONObject{
						Key: "Boolobj",
						Children: []JSONElement{
							&JSONPrimitive{
								Parent:   &JSONArray{},
								Datatype: Bool,
								Key:      "testbool",
							},
						},
					},
				},
			},
			want: "Testarray []struct{\n" +
				"Testint int `json:\"testint\"`\n" +
				"Teststring string `json:\"teststring\"`\n" +
				"Testbool bool `json:\"testbool\"`\n" +
				"} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Types Primitive",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Bool,
						Key:      "testbool",
					},
					&JSONPrimitive{
						Parent:   &JSONArray{},
						Datatype: Int,
						Key:      "testint",
					},
				},
			},
			want: "Testarray []interface{} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Types Mixed",
			fields: fields{
				Key:    "testarray",
				Parent: &JSONObject{},
				Children: []JSONElement{
					&JSONPrimitive{
						Datatype: Bool,
						Parent:   &JSONArray{},
						Key:      "testbool",
					},
					&JSONObject{
						Key: "Doesn't matter",
						Children: []JSONElement{
							&JSONPrimitive{
								Datatype: Int,
								Parent:   &JSONArray{},
								Key:      "testint",
							},
							&JSONPrimitive{
								Datatype: String,
								Parent:   &JSONArray{},
								Key:      "teststring",
							},
						},
					},
				},
			},
			want: "Testarray []interface{} `json:\"testarray\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONArray{
				Key:      tt.fields.Key,
				Parent:   tt.fields.Parent,
				Children: tt.fields.Children,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONArray.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listChildrenTypes(t *testing.T) {
	type args struct {
		c []JSONElement
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// {
		// 	name: "Two Dif Primitives",
		// 	args: args{
		// 		[]JSONElement{
		// 			&JSONPrimitive{
		// 				Datatype: String,
		// 			},
		// 			&JSONPrimitive{
		// 				Datatype: Int,
		// 			},
		// 		},
		// 	},
		// 	want: 2,
		// },
		{
			name: "Different Primitives",
			args: args{
				[]JSONElement{
					&JSONPrimitive{
						Datatype: Bool,
					},
					&JSONPrimitive{
						Datatype: String,
					},
					&JSONPrimitive{
						Datatype: Float,
					},
					&JSONPrimitive{
						Datatype: Int,
					},
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countChildrenTypes(tt.args.c)
			want := tt.want
			if !reflect.DeepEqual(got, want) {
				t.Errorf("listChildrenTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddObjects(t *testing.T) {
	arr := &JSONArray{
		Key: "Testarr",
	}
	obj := &JSONObject{
		Key: "Testobj",
	}
	bprim := &JSONPrimitive{
		Key:      "Testbool",
		Datatype: Bool,
	}
	obj.AddChild(bprim)

	arr.AddChild(obj)
	if len(arr.Children) != 1 {
		t.Error("obj has not been added")
	}
	if !reflect.DeepEqual(arr.Children[0], obj) {
		t.Error("obj has not been added the right way")
	}

	obj2 := &JSONObject{
		Key: "Testobj",
	}
	sprim := &JSONPrimitive{
		Key:      "Teststring",
		Datatype: String,
	}
	obj2.AddChild(sprim)

	// This should result in the obj being merged. More specifically, every child of
	// obj2 should be added to obj
	arr.AddChild(obj2)
	if len(arr.Children) != 1 {
		t.Error("obj2 has been added but should have been merged")
	}
	// The first child of obj should remain bprim
	if !reflect.DeepEqual(obj.Children[0], bprim) {
		t.Error("bprim was removed from obj.Children")
	}
	// The second child of obj should be sprim, which was added to obj2
	if !reflect.DeepEqual(obj.Children[1], sprim) {
		t.Error("sprim was not added to obj.Children")
	}

}
func TestAddPrimitives(t *testing.T) {
	arr := &JSONArray{
		Key: "Testarr",
	}
	sprim := &JSONPrimitive{
		Key:      "Teststring",
		Datatype: String,
	}
	arr.AddChild(sprim)

	if len(arr.Children) != 1 {
		t.Error("sprim has not been added")
	}
	if !reflect.DeepEqual(arr.Children[0], sprim) {
		t.Error("sprim has not been added the right way")
	}

	// Adding same datatype again. The child should not be added
	arr.AddChild(sprim)
	if len(arr.Children) != 1 {
		t.Error("sprim was added again but should not have been")
	}

	bprim := &JSONPrimitive{
		Key:      "Testbool",
		Datatype: Bool,
	}
	arr.AddChild(bprim)
	if len(arr.Children) != 2 {
		t.Error("brim has not been added")
	}
	if !reflect.DeepEqual(arr.Children[1], bprim) {
		t.Error("brim has not been added the right way")
	}
}
