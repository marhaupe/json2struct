package ast

import (
	"reflect"
	"testing"
)

func TestStringRootArray(t *testing.T) {
	type TestParams struct {
		name string
		got  string
		want string
	}
	tests := []func() TestParams{
		func() TestParams {
			rootArr := &JSONArray{}
			childPrimitive := &JSONPrimitive{
				Datatype: String,
			}
			rootArr.AddChild(childPrimitive)
			return TestParams{
				got:  rootArr.String(),
				want: "type JSONToStruct []string",
				name: "With Child String Primitives",
			}
		},
		func() TestParams {
			rootArr := &JSONArray{}
			childString := &JSONPrimitive{
				Datatype: String,
			}
			childInt := &JSONPrimitive{
				Datatype: Int,
			}
			rootArr.AddChild(childInt)
			rootArr.AddChild(childString)
			return TestParams{
				got:  rootArr.String(),
				want: "type JSONToStruct []interface{}",
				name: "With Child String Primitives",
			}
		},
		func() TestParams {
			rootArr := &JSONArray{}
			childArray := &JSONArray{}
			rootArr.AddChild(childArray)
			return TestParams{
				got:  rootArr.String(),
				want: "type JSONToStruct []interface{}",
				name: "With Child Array",
			}
		},
		func() TestParams {
			rootArr := &JSONArray{}
			childObj := &JSONObject{}
			nullPrim := &JSONPrimitive{
				Key:      "Testnull",
				Datatype: Null,
			}
			stringPrim := &JSONPrimitive{
				Key:      "Teststring",
				Datatype: String,
			}
			childObj.AddChild(nullPrim)
			childObj.AddChild(stringPrim)
			rootArr.AddChild(childObj)
			return TestParams{
				got: rootArr.String(),
				want: "type JSONToStruct []struct{\n" +
					"Testnull interface{} `json:\"Testnull\"`\n" +
					"Teststring string `json:\"Teststring\"`\n" +
					"}\n",
				name: "With Child Object",
			}
		},
	}
	for _, test := range tests {
		testcase := test()
		if testcase.got != testcase.want {
			t.Errorf("Test = %v Got = %v, want = %v", testcase.name, testcase.got, testcase.want)
		}
	}
}

func TestStringNodes(t *testing.T) {
	type TestParams struct {
		name string
		got  string
		want string
	}
	tests := []func() TestParams{
		func() TestParams {
			node := &JSONArray{
				Key:    "testarray",
				Parent: &JSONObject{},
			}
			childArr := &JSONArray{}
			childArr.AddChild(
				&JSONPrimitive{
					Datatype: Bool,
				},
			)
			childArr.AddChild(
				&JSONPrimitive{
					Datatype: Int,
				},
			)
			node.AddChild(childArr)
			return TestParams{
				got:  node.String(),
				name: "Array Test Arrays",
				want: "Testarray []interface{} `json:\"testarray\"`\n",
			}
		},
		func() TestParams {
			node := &JSONArray{
				Key:    "testarray",
				Parent: &JSONObject{},
			}
			obj1 := &JSONObject{}
			obj1.AddChild(&JSONPrimitive{
				Key:      "testint",
				Datatype: Int,
			})
			obj2 := &JSONObject{}
			obj2.AddChild(&JSONPrimitive{
				Key:      "teststring",
				Datatype: String,
			})
			obj3 := &JSONObject{}
			obj3.AddChild(&JSONPrimitive{
				Key:      "testbool",
				Datatype: Bool,
			})
			node.AddChild(obj1)
			node.AddChild(obj2)
			node.AddChild(obj3)
			return TestParams{
				got:  node.String(),
				name: "Different Objects",
				want: "Testarray []struct{\n" +
					"Testint int `json:\"testint\"`\n" +
					"Teststring string `json:\"teststring\"`\n" +
					"Testbool bool `json:\"testbool\"`\n" +
					"} `json:\"testarray\"`\n",
			}
		},
	}
	for _, test := range tests {
		testcase := test()
		if testcase.got != testcase.want {
			t.Errorf("Test = %v Got = %v, want = %v", testcase.name, testcase.got, testcase.want)
		}
	}
}
func TestStringPrimitives(t *testing.T) {
	type TestParams struct {
		name string
		got  string
		want string
	}
	tests := []func() TestParams{
		func() TestParams {
			arr := JSONArray{
				Key:    "testarray",
				Parent: &JSONObject{},
			}
			arr.AddChild(&JSONPrimitive{
				Datatype: Int,
			})
			arr.AddChild(&JSONPrimitive{
				Datatype: String,
			})

			return TestParams{
				name: "Different Primitives",
				got:  arr.String(),
				want: "Testarray []interface{} `json:\"testarray\"`\n",
			}
		},
		func() TestParams {
			arr := JSONArray{
				Key:    "testarray",
				Parent: &JSONObject{},
			}
			arr.AddChild(&JSONPrimitive{
				Datatype: Bool,
			})
			arr.AddChild(
				&JSONPrimitive{
					Datatype: Bool,
				})

			return TestParams{
				name: "Array Test Only Bools",
				got:  arr.String(),
				want: "Testarray []bool `json:\"testarray\"`\n",
			}
		},
	}
	for _, test := range tests {
		testcase := test()
		if testcase.got != testcase.want {
			t.Errorf("Test = %v Got = %v, want = %v", testcase.name, testcase.got, testcase.want)
		}
	}
}

func Test_listChildrenTypes(t *testing.T) {
	type args struct {
		c []Element
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Different Primitives",
			args: args{
				[]Element{
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

func TestSetGetParentArray(t *testing.T) {
	arr := &JSONArray{}
	parent := &JSONObject{}
	arr.SetParent(parent)
	if !reflect.DeepEqual(arr.GetParent(), parent) {
		t.Error("Parent has not been set or retrieved properly")
	}
}

func TestGetKeyArray(t *testing.T) {
	arr := &JSONArray{
		Key: "Testkey",
	}
	if arr.GetKey() != arr.Key {
		t.Error("Key has not been retrieved properly")
	}

}
func TestAddChild(t *testing.T) {
	type TestParams struct {
		name string
		got  Element
		want Element
	}
	tests := []func() TestParams{
		func() TestParams {
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
			return TestParams{
				name: "Add Obj To Array",
				got:  arr.Children[0],
				want: obj,
			}
		},
		func() TestParams {
			obj := &JSONObject{
				Key: "Testobj",
			}
			childObj1 := &JSONObject{
				Key: "Testobj",
			}
			childObj1.AddChild(&JSONPrimitive{
				Key:      "Teststring",
				Datatype: String,
			})
			childObj2 := &JSONObject{
				Key: "Testobj",
			}
			childObj2.AddChild(&JSONPrimitive{
				Key:      "Testint",
				Datatype: Int,
			})
			obj.AddChild(childObj1)
			obj.AddChild(childObj2)

			if len(obj.Children) != 1 {
				t.Error("childobj2 has been added but should have been merged")
			}
			got, ok := obj.Children[0].(*JSONObject)
			if !ok {
				t.Error("Casting child to object failed")
			}
			if len(got.Children) != 2 {
				t.Error("Children count is not two")
			}
			if got.Children[0].GetDatatype() != String && got.Children[0].GetKey() != "Teststring" ||
				got.Children[1].GetDatatype() != Int && got.Children[1].GetKey() != "Testint" {
				t.Error("Objects have not been merged the right way")
			}
			return TestParams{}
		},
	}
	for _, test := range tests {
		testcase := test()
		if testcase.got != testcase.want {
			t.Errorf("Test = %v Got = %v, want = %v", testcase.name, testcase.got, testcase.want)
		}
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
	if len(arr.Children) != 2 {
		t.Error("sprim was not added again but should have been")
	}

	bprim := &JSONPrimitive{
		Key:      "Testbool",
		Datatype: Bool,
	}
	arr.AddChild(bprim)
	if len(arr.Children) != 3 {
		t.Error("brim has not been added")
	}
	if !reflect.DeepEqual(arr.Children[2], bprim) {
		t.Error("brim has not been added the right way")
	}
}
