package ds

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	obj := &JSONObject{}
	wanted := "type JSONToStruct struct{\n" +
		"}\n"
	got := obj.String()
	if got != wanted {
		t.Errorf("Failed Test %v\ngot:\n%v\nwanted:\n%v", t.Name(), got, wanted)
	}

	obj = &JSONObject{
		Parent: &JSONObject{},
		Key:    "testobj",
	}
	wanted = "Testobj struct{\n" +
		"} `json:\"testobj\"`\n"
	got = obj.String()
	if got != wanted {
		t.Errorf("Failed Test %v\ngot:\n%v\nwanted:\n%v", t.Name(), got, wanted)
	}

	obj = &JSONObject{
		Parent: &JSONObject{},
		Key:    "testobj",
	}
	obj.AddChild(&JSONPrimitive{
		Key:      "teststring",
		Datatype: String,
	})
	wanted = "Testobj struct{\n" +
		"Teststring string `json:\"teststring\"`\n" +
		"} `json:\"testobj\"`\n"
	got = obj.String()
	if got != wanted {
		t.Errorf("Failed Test %v\ngot:\n%v\nwanted:\n%v", t.Name(), got, wanted)
	}
}
func TestAddChild(t *testing.T) {
	obj := &JSONObject{
		Key: "Testobj",
	}
	objToAdd := &JSONObject{
		Key: "ToAdd",
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {}
	// 		}
	// }
	obj.AddChild(objToAdd)
	if !reflect.DeepEqual(obj.Children[0], objToAdd) {
		t.Errorf("Added object to object the wrong way")
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {
	// 					"Testint": 0,
	//      }
	// 		}
	// }
	objToAdd.AddChild(&JSONPrimitive{
		Key:      "Testint",
		Datatype: Int,
	})
	if !reflect.DeepEqual(obj.Children[0], objToAdd) {
		t.Errorf("Added primitive to nested object the wrong way")
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {
	// 					"Testint": 0,
	//      },
	//			"ToAdd": {
	//			    "Testbool": false,
	//		  }
	// 		}
	// }
	secondObjToAdd := &JSONObject{
		Key: "ToAdd",
	}
	secondObjToAdd.AddChild(&JSONPrimitive{
		Key:      "Testbool",
		Datatype: Bool,
	})

	obj.AddChild(secondObjToAdd)

	mergedObj := &JSONObject{
		Parent: obj,
		Key:    "ToAdd",
	}
	mergedObj.AddChild(&JSONPrimitive{
		Key:      "Testint",
		Datatype: Int,
	})
	mergedObj.AddChild(&JSONPrimitive{
		Key:      "Testbool",
		Datatype: Bool,
	})

	if !reflect.DeepEqual(obj.Children[0], mergedObj) {
		t.Errorf("Failed to correctly merge objects\ngot:\n%v\nwanted:\n%v\n", mergedObj, obj.Children[0])
	}
}
