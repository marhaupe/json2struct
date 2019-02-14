package ast

import (
	"reflect"
	"testing"
)

func TestGetParentObjects(t *testing.T) {
	obj := &JSONObject{}
	root := &JSONObject{}
	root.AddChild(obj)
	if !reflect.DeepEqual(obj.GetParent(), root) {
		t.Error("Parent has not been set or retrieved properly")
	}
}

func TestString(t *testing.T) {
	type TestParams struct {
		name string
		got  string
		want string
	}
	tests := []func() TestParams{
		func() TestParams {
			obj := &JSONObject{
				Parent: &JSONObject{},
				Key:    "testobj",
			}
			wanted := "Testobj struct{\n" +
				"} `json:\"testobj\"`\n"
			got := obj.String()
			return TestParams{
				got:  got,
				want: wanted,
				name: "Plain Object",
			}
		},
		func() TestParams {
			obj := &JSONObject{}
			wanted := "type JSONToStruct struct{\n" +
				"}\n"
			got := obj.String()
			return TestParams{
				got:  got,
				want: wanted,
				name: "Plain Root Object",
			}
		},
		func() TestParams {
			obj := &JSONObject{
				Parent: &JSONObject{},
				Key:    "testobj",
			}
			obj.AddChild(&JSONPrimitive{
				Key:      "teststring",
				Datatype: String,
			})
			wanted := "Testobj struct{\n" +
				"Teststring string `json:\"teststring\"`\n" +
				"} `json:\"testobj\"`\n"
			got := obj.String()
			return TestParams{
				got:  got,
				want: wanted,
				name: "Object With Primitive",
			}
		},
		func() TestParams {
			obj := &JSONObject{}
			wanted := "type JSONToStruct struct{\n" +
				"}\n"
			got := obj.String()
			return TestParams{
				got:  got,
				want: wanted,
				name: "Plain Object",
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
func TestAddChildVarious(t *testing.T) {
	type TestParams struct {
		name string
		got  Element
		want Element
	}
	tests := []func() TestParams{
		func() TestParams {
			obj := &JSONObject{
				Key: "Testobj",
			}
			objToAdd := &JSONObject{
				Key: "ToAdd",
			}
			obj.AddChild(objToAdd)
			return TestParams{
				name: "Add object to object",
				want: objToAdd,
				got:  obj.Children[0],
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
